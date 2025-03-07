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

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// Routes
	http.HandleFunc("/register", handlers.RegisterHandler)
	http.HandleFunc("/login", handlers.LoginHandler)
	http.HandleFunc("/home", handlers.HomeHandler)
	http.HandleFunc("/account", handlers.AccountHandler)
	http.HandleFunc("/logout", handlers.LogoutHandler)
	http.HandleFunc("/settings", handlers.SettingsHandler)
	http.HandleFunc("/contact", handlers.ContactHandler)
	http.HandleFunc("/categories", handlers.CategoriesHandler)
	http.HandleFunc("/create-post", handlers.CreatePostHandler)
	http.HandleFunc("/", handlers.HomeHandler)

	// Démarrer le serveur
	fmt.Println("Serveur démarré sur : http://localhost:8080/home")
	log.Fatal(http.ListenAndServe(":8080", nil))
}