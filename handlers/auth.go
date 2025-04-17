package handlers

import (
	//"Forum/database"
	"database/sql"
	"fmt"
	"forum/database"
	"forum/models"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/sessions"

	//"github.com/gofrs/uuid"
	"html/template"

	"golang.org/x/crypto/bcrypt"
)

var store = sessions.NewCookieStore([]byte("super-secret-key"))

func EchecHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "templates/echec.html")
}

// on stocke la db ici une fois pour toutes
var DB *sql.DB

func SetDatabase(db *sql.DB) {
	DB = db
}

// a appeler depuis main.go pour passer la base à ce fichier

// ========================= REGISTER =========================
func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		// si on ouvre la page pour la première fois, on affiche juste le formulaire
		tmpl, _ := template.ParseFiles("templates/register.html")
		tmpl.Execute(w, nil)
		return
	}

	if r.Method == http.MethodPost {
		// on récupère les infos du formulaire
		username := r.FormValue("username")
		email := r.FormValue("email")
		password := r.FormValue("password")

		// si un des champs est vide → erreur
		if username == "" || email == "" || password == "" {
			http.ServeFile(w, r, "templates/ErrorRegister.html")
			return
		}

		// on chiffre le mot de passe pour plus de sécurité
		hashedPwd, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			http.ServeFile(w, r, "templates/ErrorRegister.html")
			return
		}

		// on enregistre le nouvel utilisateur dans la base
		db := database.GetDatabase()
		_, err = db.Exec("INSERT INTO users (username, email, password) VALUES (?, ?, ?)", username, email, string(hashedPwd))
		if err != nil {
			http.ServeFile(w, r, "templates/ErrorRegister.html")
			return
		}

		// une fois inscrit → on redirige vers la page de connexion
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}
}

// ========================= LOGIN =========================
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		tmpl, _ := template.ParseFiles("templates/login.html")
		tmpl.Execute(w, nil)
		return
	}

	if r.Method == http.MethodPost {
		email := r.FormValue("email")
		password := r.FormValue("password")

		db := database.GetDatabase()
		if db == nil {
			log.Println("❌ base non initialisée")
			http.Redirect(w, r, "/echec", http.StatusSeeOther)
			return
		}

		var id int
		var username, hashedPassword string

		err := db.QueryRow("SELECT id, username, password FROM users WHERE email = ?", email).
			Scan(&id, &username, &hashedPassword)
		if err != nil {
			log.Println("❌ Email inconnu :", err)
			http.ServeFile(w, r, "templates/ErrorLogin.html")
			return
		}

		err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
		if err != nil {
			log.Println("❌ Mot de passe incorrect :", err)
			http.ServeFile(w, r, "templates/ErrorLogin.html")
			return
		}

		// ✅ Cookie pour le nom d'utilisateur
		http.SetCookie(w, &http.Cookie{
			Name:  "username",
			Value: username,
			Path:  "/",
		})

		// ✅ Cookie pour l'ID utilisateur (sous forme de string)
		http.SetCookie(w, &http.Cookie{
			Name:  "session",
			Value: strconv.Itoa(id),
			Path:  "/",
		})

		log.Println("✅ Connexion réussie pour", username)
		http.Redirect(w, r, "/home", http.StatusSeeOther)
	}
}

// structure avec le champ FormattedDate
type PostWithFormattedDate struct {
	models.Post
	FormattedDate string
	ProfilePicture string
}

// Afficher la page d'accueil avec tous les posts
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	// on recupere le cookie username
	cookie, err := r.Cookie("username")
	username := "Invité"
	if err == nil && cookie.Value != "" {
		username = cookie.Value
	}

	// récupération des posts
	rawPosts, err := database.GetAllPosts()
	if err != nil {
		http.Redirect(w, r, "/echec", http.StatusSeeOther)
		return
	}

	// transformation avec la date formatée et la photo de profil de l'auteur
	var posts []PostWithFormattedDate
	for _, post := range rawPosts {
		var profilePicture string
		err := database.GetDatabase().QueryRow("SELECT profile_picture FROM users WHERE username = ?", post.Author).Scan(&profilePicture)
		if err != nil || profilePicture == "" {
			profilePicture = "default.jpg"
		}

		posts = append(posts, PostWithFormattedDate{
			Post:           post,
			FormattedDate:  post.Date.Format("02 Jan 2006 à 15:04"),
			ProfilePicture: profilePicture,
		})
	}

	// données envoyées au HTML
	data := struct {
		Username string
		Posts    []PostWithFormattedDate
		LoggedIn bool
	}{
		Username: username,
		Posts:    posts,
		LoggedIn: username != "Invité",
	}

	// affichage
	tmpl, err := template.ParseFiles("templates/home.html")
	if err != nil {
		log.Println("Erreur template :", err)
		http.Redirect(w, r, "/echec", http.StatusSeeOther)
		return
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		log.Println("Erreur Execute :", err)
	}
}

/*
func AccountHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session")
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	var username, email string
	err = DB.QueryRow("SELECT username, email FROM users WHERE id = ?", cookie.Value).Scan(&username, &email)
	if err != nil {
		log.Println("Erreur récupération user dans /account :", err)
		http.Redirect(w, r, "/echec", http.StatusSeeOther)
		return
	}

	data := struct {
		Username string
		Email    string
	}{
		Username: username,
		Email:    email,
	}

	tmpl, err := template.ParseFiles("templates/account.html")
	if err != nil {
		log.Println("Erreur template account :", err)
		http.Redirect(w, r, "/echec", http.StatusSeeOther)
		return
	}

	tmpl.Execute(w, data)
}*/
/*		func AccountHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session")
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	db := database.GetDatabase() // ✅ récupère la bonne connexion
	var username, email string
	err = db.QueryRow("SELECT username, email FROM users WHERE id = ?", cookie.Value).Scan(&username, &email)
	if err != nil {
		log.Println("Erreur récupération user dans /account :", err)
		http.Redirect(w, r, "/echec", http.StatusSeeOther)
		return
	}

	data := struct {
		Username string
		Email    string
	}{
		Username: username,
		Email:    email,
	}

	tmpl, err := template.ParseFiles("templates/account.html")
	if err != nil {
		log.Println("Erreur template account :", err)
		http.Redirect(w, r, "/echec", http.StatusSeeOther)
		return
	}

	tmpl.Execute(w, data)
}*/
func AccountHandler(w http.ResponseWriter, r *http.Request) {
	// Recupere le cookie de session
	cookie, err := r.Cookie("session")
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	db := database.GetDatabase()
	var username, email string
	var profilePicture sql.NullString

	// Recupere les infos de l'utilisateur
	err = db.QueryRow("SELECT username, email, profile_picture FROM users WHERE id = ?", cookie.Value).Scan(&username, &email, &profilePicture)
	if err != nil {
		log.Println("Erreur récupération user dans /account :", err)
		http.Redirect(w, r, "/echec", http.StatusSeeOther)
		return
	}

	// Si la photo est nulle, on met l'image par défaut
	pic := "default.jpg"
	if profilePicture.Valid && profilePicture.String != "" {
		pic = profilePicture.String
	}

	// Recupere le nombre de posts créés par l'utilisateur
	var postCount int
	err = db.QueryRow("SELECT COUNT(*) FROM posts WHERE author = ?", username).Scan(&postCount)
	if err != nil {
		log.Println("Erreur récupération posts dans /account :", err)
		postCount = 0
	}

	// Recupere le nombre de commentaires faits par l'utilisateur
	var commentCount int
	err = db.QueryRow("SELECT COUNT(*) FROM comments WHERE author = ?", username).Scan(&commentCount)
	if err != nil {
		log.Println("Erreur récupération commentaires dans /account :", err)
		commentCount = 0
	}

	// Recupere le nombre de likes donnés par l'utilisateur
	var likeCount int
	err = db.QueryRow("SELECT COUNT(*) FROM votes_posts WHERE user_id = ? AND vote_type = 'like'", cookie.Value).Scan(&likeCount)
	if err != nil {
		log.Println("Erreur récupération likes dans /account :", err)
		likeCount = 0
	}

	// Préparer les données pour le template
	data := struct {
		Username       string
		Email          string
		PostCount      int
		CommentCount   int
		LikeCount      int
		ProfilePicture string
	}{
		Username:       username,
		Email:          email,
		PostCount:      postCount,
		CommentCount:   commentCount,
		LikeCount:      likeCount,
		ProfilePicture: pic,
	}

	// Charger le template
	tmpl, err := template.ParseFiles("templates/account.html")
	if err != nil {
		log.Println("Erreur template account :", err)
		http.Redirect(w, r, "/echec", http.StatusSeeOther)
		return
	}

	// Execute le template avec les données
	err = tmpl.Execute(w, data)
	if err != nil {
		log.Println("Erreur Execute account.html :", err)
	}
}

// ========================= LOGOUT =========================
func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	// log de déconnexion
	log.Println("🚪 Déconnexion de l'utilisateur")

	// on supprime le cookie session
	sessionCookie := http.Cookie{
		Name:   "session",
		Value:  "",
		Path:   "/",
		MaxAge: -1, // le rend invalide
	}
	http.SetCookie(w, &sessionCookie)

	// on supprime aussi le cookie username
	usernameCookie := http.Cookie{
		Name:   "username",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	}
	http.SetCookie(w, &usernameCookie)

	// redirection vers l'accueil
	http.Redirect(w, r, "/home", http.StatusSeeOther)
}

func SettingsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		http.ServeFile(w, r, "templates/settings.html")
	}
}

func ContactHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		http.ServeFile(w, r, "templates/contact.html")
	}
}

// ========================= GET CURRENT USER =========================
func GetCurrentUser(r *http.Request) (int, string, error) {
	// on lit le cookie session
	cookie, err := r.Cookie("session")
	if err != nil {
		return 0, "", fmt.Errorf("pas de session")
	}

	// on recupere l'id depuis le cookie
	userID := cookie.Value

	// on récupère correctement la base de données
	db := database.GetDatabase()
	if db == nil {
		return 0, "", fmt.Errorf("base non initialisée")
	}

	// on cherche le pseudo dans la base
	var username string
	err = db.QueryRow("SELECT username FROM users WHERE id = ?", userID).Scan(&username)
	if err != nil {
		return 0, "", fmt.Errorf("user introuvable")
	}

	// on convertit l'id texte -> int
	id, err := strconv.Atoi(userID)
	if err != nil {
		return 0, "", fmt.Errorf("id invalide")
	}

	return id, username, nil
}

func UploadProfilePictureHandler(w http.ResponseWriter, r *http.Request) {
	// Verifie que la methode est bien POST
	if r.Method != "POST" {
		http.Redirect(w, r, "/echec.html", http.StatusSeeOther)
		return
	}

	// Recupere la session pour obtenir l'utilisateur connecte
	session, err := store.Get(r, "session-name")
	if err != nil {
		log.Println("Erreur de session :", err)
		http.Redirect(w, r, "/echec.html", http.StatusSeeOther)
		return
	}

	userID, ok := session.Values["userID"].(int)
	if !ok {
		log.Println("Utilisateur non connecte")
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Recupere le fichier depuis le formulaire
	file, handler, err := r.FormFile("profile_picture")
	if err != nil {
		log.Println("Erreur de recuperation du fichier :", err)
		http.Redirect(w, r, "/echec.html", http.StatusSeeOther)
		return
	}
	defer file.Close()

	// Cree le chemin complet pour sauvegarder l'image
	filePath := fmt.Sprintf("./static/uploads/profile_pictures/%s", handler.Filename)
	dst, err := os.Create(filePath)
	if err != nil {
		log.Println("Erreur de creation du fichier :", err)
		http.Redirect(w, r, "/echec.html", http.StatusSeeOther)
		return
	}
	defer dst.Close()

	// Copie le fichier upload dans le dossier
	if _, err := io.Copy(dst, file); err != nil {
		log.Println("Erreur de copie du fichier :", err)
		http.Redirect(w, r, "/echec.html", http.StatusSeeOther)
		return
	}

	// Met a jour la BDD avec le nom du fichier image
	db := database.GetDatabase()
	_, err = db.Exec("UPDATE users SET profile_picture = ? WHERE id = ?", handler.Filename, userID)
	if err != nil {
		log.Println("Erreur d'update BDD :", err)
		http.Redirect(w, r, "/echec.html", http.StatusSeeOther)
		return
	}

	// Redirige vers le profil ou la page precedente
	http.Redirect(w, r, "/account", http.StatusSeeOther)
}
