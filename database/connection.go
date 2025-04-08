package database

import "database/sql"

var dbInstance *sql.DB

func SetDatabase(db *sql.DB) {
	dbInstance = db
}

func GetDatabase() *sql.DB {
	return dbInstance
}