package db

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/lib/pq" // PostgreSQL driver
)

var db *sql.DB

// InitDB initializes the database connection and sets up connection pooling.
func InitDB() {
	var err error
	dbURL := os.Getenv("DATABASE_URL")
	db, err = sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}

	if err = db.Ping(); err != nil {
		log.Fatalf("Error connecting to the database: %v", err)
	}

	query := `
CREATE TABLE IF NOT EXISTS messages (
    message_id BIGINT PRIMARY KEY,
    timestamp TIMESTAMP NOT NULL,
    name TEXT,
    last_name TEXT,
    username TEXT,
    group_id BIGINT,
    user_id BIGINT,
    content TEXT
);
`
	_, err = db.Exec(query)
	if err != nil {
		log.Fatalf("Error creating messages table: %v", err)
	}
}

// GetDB returns the database connection.
func GetDB() *sql.DB {
	return db
}

// CloseDB closes the database connection.
func CloseDB() {
	if err := db.Close(); err != nil {
		log.Fatalf("Error closing the database: %v", err)
	}
	log.Println("Database connection closed.")
}
