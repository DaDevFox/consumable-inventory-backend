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

	t.Run("Put banana", func(t *testing.T) {
		putFood(FOOD{Name: "Banana", Amount: 0}, i.db)
	})
	t.Run("Put apple", func(t *testing.T) {
		putFood(FOOD{Name: "Apple", Amount: 0}, i.db)
	})
	t.Run("Put eggs", func(t *testing.T) {
		putFood(FOOD{Name: "Eggs", Amount: 0}, i.db)
	})
	t.Run("Put oranges", func(t *testing.T) {
		putFood(FOOD{Name: "Oranges", Amount: 0}, i.db)
	})
	t.Run("Put milk", func(t *testing.T) {
		putFood(FOOD{Name: "Milk", Amount: 0}, i.db)
	})
}

func (i *InventoryTestSuite) TearDownSuite() {
	i.db.Close()
}

func (i *InventoryTestSuite) TestPost() {
	t := i.T()
	postFood(FOOD{Name: "Apple", Amount: 2}, i.db)

	all := getFoods(i.db)

	for i := 0; i < len(all); i++ {
		if all[i].Name == "Apple" {
			assert.Equal(t, all[i].ID, 1)
			assert.Equal(t, all[i].Amount, 2)
		}
	}
}
