package main

import (
	"os"
	"testing"

	"github.com/codenotary/immudb/pkg/stdlib"

	secure "github.com/Dentrax/obscure-go/types"
	immudb "github.com/codenotary/immudb/pkg/client"
)

func TestStartup(t *testing.T) {
	// tests are always run with default immudb options (username immudb, password immudb, dbname defaultdb)
	immudbIP := secure.NewString(os.Args[1])

	pass := secure.NewString(os.Getenv("IMMUDB_PASSWORD"))
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

	// msg, err := Hello("Gladys")
	// if !want.MatchString(msg) || err != nil {
	// 	t.Fatalf(`Hello("Gladys") = %q, %v, want match for %#q, nil`, msg, err, want)
	// }
}
