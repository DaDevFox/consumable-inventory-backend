package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"

	secure "github.com/Dentrax/obscure-go/types"
	_ "github.com/codenotary/immudb/pkg/stdlib"
	"github.com/gin-gonic/gin"
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

func postAlbums(c *gin.Context, db *sql.DB) {
	var newFood FOOD

	// Call BindJSON to bind the received JSON to
	// newAlbum.

	err := c.BindJSON(&newFood)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Posting to food DB: " + newFood.Name + " " + strconv.Itoa(newFood.Amount) + " ")

	if newFood.Name != "" {
		// credentials.NewTLS(&tls.Config{InsecureSkipVerify: true})
		_, err := db.Exec("INSERT INTO food(name, amount) VALUES ($1, $2)", newFood.Name, newFood.Amount)
		if err != nil {
			log.Fatal(err)
		}
	}

	fmt.Println("2Posting to food DB: " + newFood.Name + " " + strconv.Itoa(newFood.Amount) + " ")

	// Add the new album to the slice.
	c.IndentedJSON(http.StatusCreated, newFood)
}

func main() {
	fmt.Println("Hello world!")

	dbName := "defaultdb"
	pass := secure.NewString("immudb") // DANGER: stored in code segment as of now AND open-source on github -- easy to find; get secure passthrough (e.g. CLI input) method to fully harden
	user := "immudb"

	connStr := secure.NewString("immudb://" + user + ":" + pass.Get() + "@127.0.0.1:3322/" + dbName + "?sslmode=disable")
	db, err := sql.Open("immudb", connStr.Get())
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	_, err = db.ExecContext(
		context.Background(),
		"CREATE TABLE IF NOT EXISTS food(id INTEGER AUTO_INCREMENT, name VARCHAR(256), amount INTEGER, PRIMARY KEY id)",
	)

	if err != nil {
		log.Fatal(err)
	}

	router := gin.Default()
	router.GET("/foods", func(c *gin.Context) {
		if err != nil {
			log.Fatal(err)
		}
		getFoods(c, db)
	})
	router.POST("/foods", func(c *gin.Context) {

		_, err := db.Query("SELECT * FROM food")
		if err != nil {
			log.Fatal(err)
		}
		postAlbums(c, db)
	})

	router.Run("localhost:5000")

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
