package db

import (
	"database/sql"

	_ "github.com/lib/pq"
)

var DB *sql.DB

// Initialize creates a connection to the PostgreSQL database
func Initialize(connStr string) error {
	var err error
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		return err
	}

	if err = DB.Ping(); err != nil {
		return err
	}
	return nil
}

// Close closes the database connection
func Close() {
	if DB != nil {
		DB.Close()
	}
}
