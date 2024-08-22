package main

import (
	"database/sql"
	"os"
	"testing"

	"github.com/codenotary/immudb/pkg/stdlib"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	secure "github.com/Dentrax/obscure-go/types"
	immudb "github.com/codenotary/immudb/pkg/client"
)

type InventoryTestSuite struct {
	suite.Suite
	db *sql.DB
}

func TestSuite(t *testing.T) {
	// run custom suite to mantain context/linearity
	suite.Run(t, new(InventoryTestSuite))
}

func (i *InventoryTestSuite) SetupSuite() {
	initAPI()

	t := i.T()

	// tests are always run with default immudb options (username immudb, password immudb, dbname defaultdb) with local immudb instances
	immudbIP := secure.NewString("127.0.0.1")

	pass := secure.NewString(os.Getenv("IMMUDB_PASSWORD"))
	user := os.Getenv("IMMUDB_USER")
	dbName := os.Getenv("IMMUDB_DB_NAME")

	// temp
	t.Log(user)
	t.Log(pass.Get())

	opts := immudb.DefaultOptions().
		WithAddress(immudbIP.Get()).
		WithPort(3322).
		WithDatabase(dbName).
		WithUsername(user).
		WithPassword(pass.Get())

	i.db = stdlib.OpenDB(opts)

	initDB(i.db)
	t.Log("DB initied and table created")
	t.Log("Inserting sample data (as subtests)")

	t.Run("Put banana", func(st *testing.T) {
		food := FOOD{Name: "Banana", Amount: 2}
		putFood(&food, i.db)
		assert.Equal(st, food.ID, 0)
	})
	t.Run("Put apple", func(st *testing.T) {
		food := FOOD{Name: "Apple", Amount: 0}
		putFood(&food, i.db)
		assert.Equal(st, food.ID, 1)
	})
	t.Run("Put eggs", func(st *testing.T) {
		food := FOOD{Name: "Eggs", Amount: 1}
		putFood(&food, i.db)
		assert.Equal(st, food.ID, 2)
	})
	t.Run("Put oranges", func(st *testing.T) {
		food := FOOD{Name: "Oranges", Amount: 0}
		putFood(&food, i.db)
		assert.Equal(st, food.ID, 3)
	})
	t.Run("Put milk", func(st *testing.T) {
		food := FOOD{Name: "Milk", Amount: 1}
		putFood(&food, i.db)
		assert.Equal(st, food.ID, 4)
	})
}

func (i *InventoryTestSuite) TearDownSuite() {
	i.db.Close()
}

func (i *InventoryTestSuite) TestPost() {
	t := i.T()
	food := FOOD{Name: "Apple", Amount: 2}
	postFood(&food, i.db)
	assert.Equal(t, food.ID, 1)

	fetched := getFood(i.db, food, false)
	assert.Equal(t, fetched.ID, 1)
	assert.Equal(t, fetched.Amount, 2)
}
