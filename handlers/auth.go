package handlers

import (
	"database/sql"
	"fmt"
	"forum/database"
	"forum/models"
	"log"
	"net/http"
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

	// transformation avec la date formatée
	var posts []PostWithFormattedDate
	for _, post := range rawPosts {
		posts = append(posts, PostWithFormattedDate{
			Post:          post,
			FormattedDate: post.Date.Format("02 Jan 2006 à 15:04"),
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
func AccountHandler(w http.ResponseWriter, r *http.Request) {
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
}
// ========================= LOGOUT =========================
func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	// on supprime le cookie en le vidant
	cookie := http.Cookie{
		Name:   "session",
		Value:  "",
		Path:   "/",
		MaxAge: -1, // ça le rend invalide
	}
	http.SetCookie(w, &cookie)

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