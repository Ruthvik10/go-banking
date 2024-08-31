package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

var testQueries *Queries
var testStore *store

var (
	dbDriver = "postgres"
	dbSource = "postgresql://root:secret@localhost:5432/bank_app?sslmode=disable"
)

func TestMain(m *testing.M) {
	db, err := sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatalf("cannot connect to the database: %s", err.Error())
	}
	testQueries = New(db)
	testStore = NewStore(db)
	os.Exit(m.Run())
}
