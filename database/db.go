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

	// Connexion initiale à MySQL sans spécifier la base de données
	dsn := "root:root_password@tcp(localhost:3306)/" // Remplace "root_password" par ton mot de passe MySQL
	tempDB, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("Erreur de connexion à MySQL:", err)
	}
	defer tempDB.Close()

	// Création de la base de données si elle n'existe pas
	_, err = tempDB.Exec("CREATE DATABASE IF NOT EXISTS mydb")
	if err != nil {
		log.Fatal("Erreur lors de la création de la base de données:", err)
	}

	fmt.Println("✅ Base de données 'mydb' prête !")

	// Connexion finale avec la base de données 'mydb'
	dsn = "root:root_password@tcp(localhost:3306)/mydb"
	DB, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("Erreur de connexion à MySQL:", err)
	}

	// Vérifier la connexion
	if err := DB.Ping(); err != nil {
		log.Fatal("Impossible de se connecter à MySQL:", err)
	}

	fmt.Println("✅ Connexion réussie à MySQL !")

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

// // Fonction pour ajouter un utilisateur dans la base de données
// func AddUser(id, username, email, password string) {
// 	// Requête SQL d'insertion
// 	query := `INSERT INTO users (id, username, email, password) VALUES (?, ?, ?, ?)`

// 	// Exécution de la requête d'insertion
// 	_, err := DB.Exec(query, id, username, email, password)
// 	if err != nil {
// 		log.Fatal("Erreur lors de l'insertion de l'utilisateur:", err)
// 	}
// 	fmt.Println("Utilisateur ajouté avec succès!")
// }

// func AddUser(id, username, email, password string) {
//     fmt.Println("📌 Tentative d'ajout de l'utilisateur :", username, email)

//     query := `INSERT INTO users (id, username, email, password) VALUES (?, ?, ?, ?)`

//     _, err := DB.Exec(query, id, username, email, password)
//     if err != nil {
//         log.Fatal("Erreur lors de l'insertion de l'utilisateur:", err)
//     }
//     fmt.Println("✅ Utilisateur ajouté avec succès!")
// }
