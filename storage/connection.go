package storage

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

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
	// Note: Saving credentials in environment variables is convenient, but not
	// secure - consider a more secure solution such as
	// Cloud Secret Manager (https://cloud.google.com/secret-manager) to help
	// keep secrets safe.
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

	configureConnectionPool(dbPool)

	return dbPool, nil
}

func configureConnectionPool(db *sql.DB) {
	// [START cloud_sql_postgres_databasesql_limit]
	// Set maximum number of connections in idle connection pool.
	db.SetMaxIdleConns(5)

	// Set maximum number of open connections to the database.
	db.SetMaxOpenConns(7)
	// [END cloud_sql_postgres_databasesql_limit]

	// [START cloud_sql_postgres_databasesql_lifetime]
	// Set Maximum time (in seconds) that a connection can remain open.
	db.SetConnMaxLifetime(1800 * time.Second)
	// [END cloud_sql_postgres_databasesql_lifetime]

	// [START cloud_sql_postgres_databasesql_backoff]
	// database/sql does not support specifying backoff
	// [END cloud_sql_postgres_databasesql_backoff]
	// [START cloud_sql_postgres_databasesql_timeout]
	// The database/sql package currently doesn't offer any functionality to
	// configure connection timeout.
	// [END cloud_sql_postgres_databasesql_timeout]
}
