package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

// InitDB initialise la base de données MySQL
func InitDB() {
	var err error

	// Connexion à MySQL (modifie les identifiants selon ton setup)
	dsn := "root:root@tcp(localhost:3306)/mydb"
	DB, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("Erreur de connexion à MySQL:", err)
	}

	// Vérifier la connexion
	if err := DB.Ping(); err != nil {
		log.Fatal("Impossible de se connecter à MySQL:", err)
	}

	fmt.Println("Connexion réussie à MySQL !")

	// Création de la table users si elle n'existe pas
	statement, err := DB.Prepare(`
		CREATE TABLE IF NOT EXISTS users (
			id VARCHAR(36) PRIMARY KEY,
			username VARCHAR(255) NOT NULL,
			email VARCHAR(255) UNIQUE NOT NULL,
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
