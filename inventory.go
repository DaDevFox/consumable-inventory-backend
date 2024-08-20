package main

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"strconv"

	secure "github.com/Dentrax/obscure-go/types"
	immudb "github.com/codenotary/immudb/pkg/client"
	"github.com/codenotary/immudb/pkg/stdlib"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

// use:
// https://github.com/codenotary/immudb-client-examples/tree/master/go/todos-sample-stdlib
// https://github.com/codenotary/immudb-client-examples/tree/master/go
// https://play.codenotary.com/?topic=cli%2Fselect&live=true
// getAlbums responds with the list of all albums as JSON.

type FOOD struct {
	ID     int
	Name   string
	Amount int
}

func getFoods(c *gin.Context, db *sql.DB) {
	var id int
	var name string
	var amount int

	var foods []FOOD

	rows, err := db.Query("SELECT id, name, amount FROM food")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&id, &name, &amount)
		if err != nil {
			log.Fatal(err)
		}

		foods = append(foods, FOOD{ID: id, Name: name, Amount: amount})
	}

	c.IndentedJSON(http.StatusOK, foods)
}

func postFood(c *gin.Context, db *sql.DB) {
	var newFood FOOD

	err := c.BindJSON(&newFood)
	if err != nil {
		c.AbortWithStatus(400)
		log.Info("Error while parsing JSON; malformed request")
		log.Error(err)
		return
	}

	if newFood.Name != "" {
		_, err := db.Exec("INSERT INTO food(name, amount) VALUES ($1, $2)", newFood.Name, newFood.Amount)
		if err != nil {
			log.Fatal(err)
		}
	}

	fmt.Println("Posting to food DB: " + newFood.Name + " " + strconv.Itoa(newFood.Amount) + " ")

	c.IndentedJSON(http.StatusCreated, newFood)
}

func printHelp() {
	log.Error("Invalid;\n\tUsage: <immudb server ip address:required> <ip address to run on: required> [port]")
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
		"CREATE UNIQUE INDEX IF NOT EXISTS ON food(id);",
	)

	_, err = db.ExecContext(
		context.Background(),
		"CREATE INDEX IF NOT EXISTS ON food(name);",
	)

	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	log.SetFormatter(&log.TextFormatter{})
	log.Debug("Init server interface API")

	if len(os.Args) <= 2 {
		printHelp()
		return
	} //TODO: implement port cli arg

	immudbIP := secure.NewString(os.Args[1])
	selfIP := secure.NewString(os.Args[2])

	pass := secure.NewString("immudb") // DANGER: stored in code segment as of now AND open-source on github -- easy to find; get secure passthrough (e.g. CLI input) method to fully harden
	user := "immudb"
	dbName := "defaultdb"

	opts := immudb.DefaultOptions().
		WithAddress(immudbIP.Get()).
		WithPort(3322). // TODO: implement immudb port cli arg
		WithDatabase(dbName).
		WithUsername(user).
		WithPassword(pass.Get())

	db := stdlib.OpenDB(opts)
	defer db.Close()

	initDB(db)

	router := gin.Default()
	router.GET("/foods", func(c *gin.Context) {
		getFoods(c, db)
	})
	router.POST("/foods", func(c *gin.Context) {

		_, err := db.Query("SELECT * FROM food")
		if err != nil {
			log.Fatal(err)
		}
		postFood(c, db)
	})

	router.Run(selfIP.Get() + ":5000")

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
	// 	log.Fatal(err)
	// }
	// fmt.Printf("Sucessfully got verified entry: ('%s', '%s') @ tx %d\n", entry.Key, entry.Value, entry.Tx)

}
