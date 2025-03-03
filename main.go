package main

import (
	"fmt"
	"log"
	"net/http"

	"forum/database"
	"forum/handlers"
)

func main() {
	// Initialisation de la BDD
	database.InitDB()

	// Routes
	http.HandleFunc("/register", handlers.RegisterHandler)
	http.HandleFunc("/login", handlers.LoginHandler)
	http.HandleFunc("/home", handlers.HomeHandler)
	http.HandleFunc("/account", handlers.AccountHandler)
	http.HandleFunc("/logout", handlers.LogoutHandler)
	http.HandleFunc("/settings", handlers.SettingsHandler)

	// Démarrer le serveur
	fmt.Println("Serveur démarré sur : http://localhost:8080/home")
	log.Fatal(http.ListenAndServe(":8080", nil))
}