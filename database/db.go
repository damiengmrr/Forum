package database

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

// initDB initialise la base de données SQLite
func InitDB() {
	var err error
	DB, err = sql.Open("sqlite3", "./forum.db")
	if err != nil {
		log.Fatal(err)
	}

	// Création de la table users si elle n'existe pas
	statement, err := DB.Prepare(`
		CREATE TABLE IF NOT EXISTS users (
			id TEXT PRIMARY KEY,
			username TEXT NOT NULL,
			email TEXT UNIQUE NOT NULL,
			password TEXT NOT NULL
		)
	`)
	if err != nil {
		log.Fatal(err)
	}
	statement.Exec()
}
