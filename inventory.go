package main

import (
	"context"
	"database/sql"
	"net/http"
	"os"
	"strconv"

	secure "github.com/Dentrax/obscure-go/types"
	immudb "github.com/codenotary/immudb/pkg/client"
	"github.com/codenotary/immudb/pkg/stdlib"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
)

// use:
// https://github.com/codenotary/immudb-client-examples/tree/master/go/todos-sample-stdlib
// https://github.com/codenotary/immudb-client-examples/tree/master/go
// https://play.codenotary.com/?topic=cli%2Fselect&live=true
// getAlbums responds with the list of all albums as JSON.
//
// EVENTUALLY:
// https://computingpost.medium.com/run-immudb-sql-and-key-value-database-on-docker-kubernetes-15f22391dca5

type FOOD struct {
	ID     int
	Name   string
	Amount int
}

func getFoods(db *sql.DB) []FOOD {
	var id int
	var name string
	var amount int

	var foods []FOOD

	rows, err := db.Query("SELECT id, name, amount FROM food")
	if err != nil {
		log.Error(err)
		return []FOOD{}
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&id, &name, &amount)
		if err != nil {
			log.Error(err)
			continue
		}

		foods = append(foods, FOOD{ID: id, Name: name, Amount: amount})
	}

	return foods
}

func getFood(db *sql.DB, food FOOD, requireIDMatch bool) *FOOD {
	log.Debug("Query: " + `SELECT id, name, amount FROM food WHERE name='` + food.Name + `'`)
	res, err := db.Query(`SELECT id, name, amount FROM food WHERE name='` + food.Name + `'`)
	if err != nil {
		log.Error(err)
		return nil
	}

	if !res.Next() {
		return nil
	}

	var id int
	var name string
	var amount int
	err = res.Scan(&id, &name, &amount)
	if err != nil {
		log.Error(err)
		return nil
	}

	log.Debug("ID scan result: name=" + name + "; id=" + strconv.Itoa(id))

	// also check unique if requireIDMatch
	if requireIDMatch && food.ID != id {
		return nil
	} else {
		return &FOOD{ID: id, Name: name, Amount: amount}
	}
}

func foodExists(db *sql.DB, food FOOD, requireIDMatch bool) bool {
	return getFood(db, food, requireIDMatch) != nil
}

func postFood(food *FOOD, db *sql.DB) {
	current := getFood(db, *food, false)
	// TODO: allow creating with POST (same functionality as PUT)
	// ----
	// see https://stackoverflow.com/questions/630453/what-is-the-difference-between-post-and-put-in-http
	// ----
	// POST to a URL creates a child resource at a server defined URL.
	// PUT to a URL creates/replaces the resource in its entirety at the client defined URL.
	// PATCH to a URL updates part of the resource at that client defined URL.

	// confirm item name exists (maybe also req. knowning id?)
	if current == nil {
		log.Info("POST request made for resource which did not exist: id=" + strconv.Itoa(food.ID) + ", name=" + food.Name)
		panic("Not found")
	}

	// pick up ID from current -- user may not know it; only req. name as identification for now
	log.Debug("POSTing name=" + current.Name + " id=" + strconv.Itoa(current.ID))
	food.ID = current.ID

	// update the item
	_, err := db.Exec("UPSERT INTO food (id, name, amount) VALUES ($1, $2, $3)", current.ID, current.Name, food.Amount)

	// update return (food) to accurately reflect changes
	food = current
	if err != nil {
		log.Info("Error posting to DB; post may have been malformed")
		log.Error(err)
	}
}

// requires struct with everything except ID set; ID is set during execution to assigned ID if successful
func putFood(newFood *FOOD, db *sql.DB) {

	log.Debug("PUTting into food DB: " + newFood.Name + " " + strconv.Itoa(newFood.Amount))

	if foodExists(db, *newFood, false) {
		log.Info("PUT request made for resource which already exists: name=" + newFood.Name)
		panic("Resource exists")
	}

	if newFood.Name != "" {
		_, err := db.Exec("INSERT INTO food(name, amount) VALUES ($1, $2)", newFood.Name, newFood.Amount)
		if err != nil {
			log.Info("Error putting in DB; put may have been malformed")
			log.Error(err)
		}

		newFood.ID = getFood(db, *newFood, false).ID
	}
}

func printHelp() {
	log.Error("Invalid;\n\tUsage: <immudb server ip address:required> <ip address to run on: required> [port]")
}

func initAPI() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
	}

	log.SetFormatter(&log.TextFormatter{})
	log.Debug("Init server interface API")
}

func initDB(db *sql.DB) {
	_, err := db.ExecContext(
		context.Background(),
		"CREATE TABLE IF NOT EXISTS food(id INTEGER AUTO_INCREMENT, name VARCHAR(256), amount INTEGER, PRIMARY KEY id)",
	)

	if err != nil {
		log.Fatal(err)
	}

	_, err = db.ExecContext(
		context.Background(),
		"CREATE UNIQUE INDEX IF NOT EXISTS ON food(id);", // TODO: determine if this is acc necessary
	)

	_, err = db.ExecContext(
		context.Background(),
		"CREATE UNIQUE INDEX IF NOT EXISTS ON food(name);",
	)

	if err != nil {
		log.Error(err)
	}
}

func main() {
	if len(os.Args) <= 2 {
		printHelp()
		return
	} //TODO: implement port cli args

	initAPI()

	immudbIP := secure.NewString(os.Args[1])
	selfIP := secure.NewString(os.Args[2])

	// TODO: close the loop so that this protection actually matters (protected in memory now but acc source is right there on disk too -- need encrypted store)
	pass := secure.NewString(os.Getenv("IMMUDB_PASSWORD")) // DANGER: stored in code segment as of now AND open-source on github -- easy to find; get secure passthrough (e.g. CLI input) method to fully harden
	user := os.Getenv("IMMUDB_USER")
	dbName := os.Getenv("IMMUDB_DB_NAME")

	opts := immudb.DefaultOptions().
		WithAddress(immudbIP.Get()).
		WithPort(3322).
		WithDatabase(dbName).
		WithUsername(user).
		WithPassword(pass.Get())

	db := stdlib.OpenDB(opts)
	defer db.Close()

	initDB(db)

	router := gin.Default()
	router.GET("/foods", func(c *gin.Context) {
		// always returns OK; no critical failure points
		foods := getFoods(db)
		c.IndentedJSON(http.StatusOK, foods)
	})
	router.PUT("/foods", func(c *gin.Context) {
		var food FOOD
		defer func() {
			if r := recover(); r != nil {
				log.Error("Panic!")
				if r == "Bad request" {
					c.AbortWithStatus(http.StatusBadRequest)
				}
				if r == "Resource exists" {
					c.AbortWithStatus(http.StatusConflict)
				}
			} else {
				c.IndentedJSON(http.StatusOK, food)
			}
		}()

		err := c.BindJSON(&food)
		if err != nil {
			log.Info("Error while parsing JSON; malformed request")
			log.Error(err)
			panic("Bad request")
		}
		putFood(&food, db)
	})
	router.POST("/foods", func(c *gin.Context) {
		var food FOOD
		defer func() {
			if r := recover(); r != nil {
				log.Error("Panic!")
				if r == "Bad request" {
					c.AbortWithStatus(http.StatusBadRequest)
				}
				if r == "Not found" {
					c.AbortWithStatus(http.StatusNotFound)
				}
			} else {
				c.IndentedJSON(http.StatusOK, food)
			}
		}()

		err := c.BindJSON(&food)
		if err != nil {
			log.Info("Error while parsing JSON; malformed request")
			log.Error(err)
			panic("Bad request")
		}
		postFood(&food, db)
	})

	// TODO: set up PKI and make this RunTLS to use https
	router.Run(selfIP.Get() + ":5000")
}

// alternate:
// opts := immudb.DefaultOptions().WithAddress("127.0.0.1").WithPort(3322)
// client := immudb.NewClient().WithOptions(opts)
//
// err := client.OpenSession(context.Background(), []byte(user), []byte(pass.Get()), "defaultdb")
// if err != nil {
// 	log.Fatal(err)
// }
// // ensure connection is closed
// defer client.CloseSession(context.Background())
//
// // write an entry
// // upon submission, the SDK validates proofs and updates the local state under the hood
// hdr, err := client.VerifiedSet(context.Background(), []byte(key), []byte("immutable world"))
// if err != nil {
// 	log.Fatal(err)
// }
// fmt.Printf("Sucessfully set a verified entry: ('%s', '%s') @ tx %d\n", []byte(key), []byte("immutable world"), hdr.Id)
//
// // read an entry
// // upon submission, the SDK validates proofs and updates the local state under the hood
// entry, err := client.VerifiedGet(context.Background(), []byte(key))
// if err != nil {
