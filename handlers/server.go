package handlers

import (
	"database/sql"
	"fmt"
	"forum/database"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

func StartServer() {
	// ouverture de la base
	db, err := sql.Open("sqlite3", "./forum.db")
	if err != nil {
		log.Fatal("Erreur de connexion à la base :", err)
	}

	// Connexion globale pour tout le projet
	database.SetDatabase(db)
	log.Println("✅ Connexion à forum.db établie")

	// DEBUG : affichage utilisateurs
	rows, err := db.Query("SELECT id, username, email FROM users")
	if err != nil {
		log.Println("❌ Erreur SELECT users au démarrage :", err)
	} else {
		defer rows.Close()
		for rows.Next() {
			var id int
			var u, e string
			rows.Scan(&id, &u, &e)
			log.Println("👤 Utilisateur trouvé :", id, u, e)
		}
	}

	// fichiers statiques
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// routes
	http.HandleFunc("/register", RegisterHandler)
	http.HandleFunc("/login", LoginHandler)
	http.HandleFunc("/home", HomeHandler)
	http.HandleFunc("/account", AccountHandler)
	http.HandleFunc("/logout", LogoutHandler)
	http.HandleFunc("/settings", SettingsHandler)
	http.HandleFunc("/contact", ContactHandler)
	http.HandleFunc("/categories", CategoriesHandler)
	http.HandleFunc("/create-post", CreatePostHandler)
	http.HandleFunc("/posts", TimeHandlers)
	http.HandleFunc("/post", PostHandler)
	http.HandleFunc("/echec", EchecHandler)
	http.HandleFunc("/submit-post", CreatePostHandler)
	http.HandleFunc("/comment/reply", CommentReplyHandler)
	http.HandleFunc("/comment/like", CommentLikeHandler)
	http.HandleFunc("/comment/dislike", CommentDislikeHandler)
	http.HandleFunc("/like", LikeHandler)
	http.HandleFunc("/dislike", DislikeHandler)
	http.HandleFunc("/test-sessions", TestSessionHandler)
	http.HandleFunc("/edit-profile", EditProfileHandler)
	http.HandleFunc("/change-password", ChangePasswordHandler)
	http.HandleFunc("/delete-post", DeletePostHandler)
	http.HandleFunc("/", EchecHandler)

	// lancement serveur HTTPS
	fmt.Println("============================================")
	fmt.Println("🚀 Lancement du serveur FORUM")
	fmt.Println("🌐 Adresse : https://localhost:8443/home")
	fmt.Println("✅ Statut  : EN LIGNE (HTTPS sécurisé)")
	fmt.Println("📌 Pour arrêter : Ctrl + C")
	fmt.Println("============================================")

	certFile := "cert.pem" // Assure-toi d'avoir ton fichier de certificat
	keyFile := "key.pem"   // Et ta clé privée à la racine du projet

	log.Fatal(http.ListenAndServeTLS(":8443", certFile, keyFile, nil))
}
