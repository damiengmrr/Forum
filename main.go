package main

import (
	"fmt"
	"log"
	"net/http"

	//"forum/database"
	"forum/handlers"
)

func main() {
	// Initialisation de la BDD
	//database.InitDB()
	//database.AddUser("1", "testuser", "test@example.com", "password123")

	//database.InitDB()

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
	http.HandleFunc("/posts", handlers.TimeHandlers)
	http.HandleFunc("/post/", handlers.PostHandler)
	http.HandleFunc("/", handlers.EchecHandler)
	http.HandleFunc("/echec", handlers.EchecHandler)
	http.HandleFunc("/submit-post", handlers.CreatePostHandler)


	fmt.Println("Serveur démarré sur http://localhost:8080/home")

	// Démarrage du serveur HTTP
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("Erreur lors du démarrage du serveur:", err)
	}
}
