package storage

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

const (
	host     = "35.228.99.104"
	port     = 5432
	user     = "postgres"
	password = "admin"
	dbname   = "postgres"
)

func Connect() (*sql.DB, error) {
	var db *sql.DB
	var err error

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	fmt.Println("Сonnected to database.")

	return db, nil
}
