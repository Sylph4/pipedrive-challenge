package storage

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/v4/stdlib"
)

func Connect() (*sql.DB, error) {
	mustGetenv := func(k string) string {
		v := os.Getenv(k)
		if v == "" {
			log.Fatalf("Warning: %s environment variable not set.\n", k)
		}
		return v
	}

	var (
		dbUser         = mustGetenv("DB_USER")              // e.g. 'my-db-user'
		dbPwd          = mustGetenv("DB_PASS")              // e.g. 'my-db-password'
		unixSocketPath = mustGetenv("INSTANCE_UNIX_SOCKET") // e.g. '/cloudsql/project:region:instance'
		dbName         = mustGetenv("DB_NAME")              // e.g. 'my-database'
	)

	dbURI := fmt.Sprintf("user=%s password=%s database=%s host=%s",
		dbUser, dbPwd, dbName, unixSocketPath)

	// dbPool is the pool of database connections.
	dbPool, err := sql.Open("pgx", dbURI)
	if err != nil {
		return nil, fmt.Errorf("sql.Open: %v", err)
	}

	m, err := migrate.New(
		"file://migrations",
		"postgres://"+dbUser+":"+dbPwd+"@"+unixSocketPath[1:]+"/"+dbName+"?sslmode=disable")
	if err != nil {
		fmt.Println(err, "1")
	}
	if err := m.Up(); err != nil {
		fmt.Println(err, "2")
	}

	configureConnectionPool(dbPool)

	return dbPool, nil
}

func configureConnectionPool(db *sql.DB) {
	db.SetMaxIdleConns(5)

	db.SetMaxOpenConns(7)

	db.SetConnMaxLifetime(1800 * time.Second)
}
