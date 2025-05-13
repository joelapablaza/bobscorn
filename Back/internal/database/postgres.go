package database

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

const (
	maxOpenConns    = 25
	maxIdleConns    = 25
	connMaxLifetime = 5 * time.Minute
)

func NewPostgreSQL(connStr string) (*sql.DB, error) {
	db, err := sql.Open("pgx", connStr)
	if err != nil {
		return nil, fmt.Errorf("error opening database connection: %w", err)
	}

	db.SetMaxOpenConns(maxOpenConns)
	db.SetMaxIdleConns(maxIdleConns)
	db.SetConnMaxLifetime(connMaxLifetime)

	// Ping the database to verify connection
	if err = db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("error pinging database: %w", err)
	}

	log.Println("Successfully connected to PostgreSQL database.")
	return db, nil
}

func InitializeSchema(db *sql.DB) error {
	createTableQuery := `
    CREATE TABLE IF NOT EXISTS client_requests (
        client_ip TEXT PRIMARY KEY,
        last_request_at TIMESTAMP WITH TIME ZONE NOT NULL
    );`

	_, err := db.Exec(createTableQuery)
	if err != nil {
		return fmt.Errorf("error creating client_requests table: %w", err)
	}

	createIndexQuery := `
    CREATE INDEX IF NOT EXISTS idx_client_requests_last_request_at 
    ON client_requests (last_request_at);`

	_, err = db.Exec(createIndexQuery)
	if err != nil {
		return fmt.Errorf("error creating index on client_requests(last_request_at): %w", err)
	}

	return nil
}
