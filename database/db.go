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
    			username VARCHAR(100),
    			email VARCHAR(100) UNIQUE,
   				password VARCHAR(255)

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

func AddUser(id, username, email, hashedPassword string) error {
	query := `INSERT INTO users (id, username, email, password) VALUES (?, ?, ?, ?)`

	_, err := DB.Exec(query, id, username, email, hashedPassword)
	if err != nil {
		log.Println("❌ Erreur lors de l'insertion de l'utilisateur :", err)
		return err
	}

	// Log de succès
	log.Printf("✅ Utilisateur ajouté : %s (%s)", username, email)
	return nil
}

