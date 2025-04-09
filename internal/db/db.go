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

    log.Println("Database connection established successfully.")
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