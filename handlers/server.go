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
		log.Fatal("❌ Erreur ouverture base :", err)
	}

	database.SetDatabase(db)
	log.Println("✅ Connexion à forum.db établie")

	// verification connexion BDD
	if !database.CheckDatabaseConnection() {
		log.Fatal("❌ Connexion à la base de données échouée. Arrêt serveur.")
	}

	// affichage utilisateurs existants
	showExistingUsers(db)

	// gestion des fichiers statiques
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// routes principales 🧭
	http.HandleFunc("/", EchecHandler)
	http.HandleFunc("/home", HomeHandler)
	http.HandleFunc("/account", AccountHandler)
	http.HandleFunc("/register", RegisterHandler)
	http.HandleFunc("/login", LoginHandler)
	http.HandleFunc("/logout", LogoutHandler)

	// routes posts et commentaires 📝
	http.HandleFunc("/create-post", CreatePostHandler)
	http.HandleFunc("/submit-post", CreatePostHandler) // doublon ?
	http.HandleFunc("/post", PostHandler)
	http.HandleFunc("/posts", TimeHandlers)
	http.HandleFunc("/categories", CategoriesHandler)
	http.HandleFunc("/comment/reply", CommentReplyHandler)
	http.HandleFunc("/comment/like", CommentLikeHandler)
	http.HandleFunc("/comment/dislike", CommentDislikeHandler)

	// routes likes / dislikes ❤️
	http.HandleFunc("/like", LikeHandler)
	http.HandleFunc("/dislike", DislikeHandler)

	// outils et pages annexes 🧩
	http.HandleFunc("/settings", SettingsHandler)
	http.HandleFunc("/contact", ContactHandler)
	http.HandleFunc("/test-sessions", TestSessionHandler)
	http.HandleFunc("/edit-profile", EditProfileHandler)
	http.HandleFunc("/change-password", ChangePasswordHandler)
	http.HandleFunc("/delete-post", DeletePostHandler)
	http.HandleFunc("/upload-profile-picture", UploadProfilePictureHandler)
	http.HandleFunc("/", EchecHandler)

	// lancement du serveur 🧩
	printServerStart()
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// fonction pour afficher les users au démarrage
func showExistingUsers(db *sql.DB) {
	rows, err := db.Query("SELECT id, username, email FROM users")
	if err != nil {
		log.Println("❌ Erreur SELECT users au démarrage :", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var username, email string
		err := rows.Scan(&id, &username, &email)
		if err != nil {
			log.Println("❌ Erreur scan user :", err)
			continue
		}
		log.Printf("👤 Utilisateur : ID %d | %s | %s\n", id, username, email)
	}
}

// affichage serveur démarré
func printServerStart() {
	fmt.Println("============================================")
	fmt.Println("🚀 Lancement du serveur FORUM")
	fmt.Println("🌐 Adresse : http://localhost:8080/home")
	fmt.Println("✅ Statut  : EN LIGNE")
	fmt.Println("📌 Pour arrêter : Ctrl + C")
	fmt.Println("============================================")
}
