package database

import (
	"database/sql"
	"iron-stream/internal/utils"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func ConnectDB() {
	var err error

	dbPath := utils.GetEnv("DB_PATH")
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
