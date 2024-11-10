package store

import (
	"log"
	"os"
	"testing"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

const (
	dsn = "postgresql://root:secret@localhost:5432/bank_app?sslmode=disable"
)

var store *Store

func TestMain(m *testing.M) {
	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		log.Fatalln(err)
	}
	store = NewStore(db)
	os.Exit(m.Run())
}
