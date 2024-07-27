package database

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func ConnectDB(path string) {
	var err error

	dbPath := os.Getenv(path)
	if dbPath == "" {
		log.Fatal("DB_PATH environment variable is not set")
	}

	DB, err = sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatalf("Failed to open the database at %s: %v", dbPath, err)
	}

	if err = DB.Ping(); err != nil {
		log.Fatalf("Failed to connect to the database at %s: %v", dbPath, err)
	}
}
