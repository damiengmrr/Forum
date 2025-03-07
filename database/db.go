package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

// InitDB initialise la base de données SQLite3
func InitDB() {
	var err error

	// Connexion à SQLite3 (modifie le chemin vers le fichier de base de données)
	dsn := "file:mydb.sqlite?cache=shared&mode=rwc" // "root:root_password@tcp(localhost:3306)/mydb"
	DB, err = sql.Open("sqlite3", dsn)
	if err != nil {
		log.Fatal("Erreur de connexion à SQLite3:", err)
	}

	// Vérifier la connexion
	if err := DB.Ping(); err != nil {
		log.Fatal("Impossible de se connecter à SQLite3:", err)
	}

	fmt.Println("Connexion réussie à SQLite3 !")

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
		log.Fatal("Erreur de création de table:", err)
	}

	_, err = statement.Exec()
	if err != nil {
		log.Fatal("Erreur d'exécution de la requête:", err)
	}

	fmt.Println("✅ Table 'users' prête !")
}
