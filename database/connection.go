package database

import (
	"database/sql"
	"log"
)

var dbInstance *sql.DB

// set la connexion globale a la bdd
func SetDatabase(db *sql.DB) {
	if db == nil {
		log.Println("❌ tentative de set une db nulle")
		return
	}
	dbInstance = db
	log.Println("✅ connexion a la base de donnees initialisee avec succes")
}

// recupere la connexion globale
func GetDatabase() *sql.DB {
	if dbInstance == nil {
		log.Println("⚠️ la base de donnees n'est pas initialisee !")
	}
	return dbInstance
}

// check si la bdd est bien connectee
func CheckDatabaseConnection() bool {
	if dbInstance == nil {
		log.Println("❌ aucune connexion a la base de donnees")
		return false
	}
	if err := dbInstance.Ping(); err != nil {
		log.Println("❌ ping bdd echoue:", err)
		return false
	}
	log.Println("✅ la connexion a la base de donnees est operationnelle")
	return true
}
