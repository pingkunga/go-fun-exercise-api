package postgres

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/lib/pq"
)

type Postgres struct {
	Db *sql.DB
}

func New() (*Postgres, error) {
	//replace connection string with ENV

	databaseSource := os.Getenv("DB_CONN")
	db, err := sql.Open("postgres", databaseSource)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	return &Postgres{Db: db}, nil
}
