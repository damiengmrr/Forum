package handlers

import (
	"database/sql"
	"fmt"
	"forum/models"
	"log"
	"net/http"
	"strconv"
	"strings"

	//"github.com/gofrs/uuid"
	"html/template"

	"golang.org/x/crypto/bcrypt"
)

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
		tmpl, _ := template.ParseFiles("templates/register.html")
		tmpl.Execute(w, nil)
		return
	}

	if r.Method == http.MethodPost {
		username := strings.TrimSpace(r.FormValue("username"))
		email := strings.TrimSpace(r.FormValue("email"))
		password := r.FormValue("password")

		log.Println("Tentative inscription avec :", username, email)

		if username == "" || email == "" || password == "" {
			log.Println("Champs vides dans le formulaire")
			http.ServeFile(w, r, "templates/ErrorRegister.html")
			return
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			log.Println("Erreur hash mot de passe:", err)
			http.ServeFile(w, r, "templates/ErrorRegister.html")
			return
		}

		res, err := DB.Exec("INSERT INTO users (username, email, password) VALUES (?, ?, ?)", username, email, string(hashedPassword))
		if err != nil {
			log.Println("Erreur SQL INSERT Register:", err)
			http.ServeFile(w, r, "templates/ErrorRegister.html")
			return
		}

		id, _ := res.LastInsertId()
		log.Println("✅ User enregistré avec succès, ID =", id)
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

		var id string
		var username, hashedPassword string

		// on récupère le user
		err := DB.QueryRow("SELECT id, username, password FROM users WHERE email = ?", email).Scan(&id, &username, &hashedPassword)
		if err != nil {
			log.Println("Email inconnu :", err)
			http.ServeFile(w, r, "templates/ErrorLogin.html")
			return
		}

		// vérifie le mot de passe
		err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
		if err != nil {
			log.Println("Mot de passe incorrect :", err)
			http.ServeFile(w, r, "templates/ErrorLogin.html")
			return
		}

		// crée le cookie de session
		cookie := http.Cookie{
			Name:  "session",
			Value: id,
			Path:  "/",
		}
		http.SetCookie(w, &cookie)

		log.Println("Connexion réussie pour", username)
		http.Redirect(w, r, "/home", http.StatusSeeOther)
	}
}

// Afficher la page d'accueil avec tous les posts
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	id, username, err := GetCurrentUser(r)

	if err != nil {
		log.Println("Invité")
	} else {
		log.Println("Connecté :", username, "ID :", id)
	}

	// ici tu prepares les donnees a envoyer au template
	data := struct {
		Posts    []models.Post
		Username string
		LoggedIn bool
	}{
		Posts:    posts,
		Username: username,
		LoggedIn: err == nil,
	}

	// on affiche la page d'accueil
	tmpl, err := template.ParseFiles("templates/home.html")
	if err != nil {
		http.Error(w, "erreur de template", http.StatusInternalServerError)
		fmt.Print(err)
		return
	}
	tmpl.Execute(w, data)
}

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

func CategoriesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		http.ServeFile(w, r, "templates/categories.html")
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

	// on cherche le pseudo depuis la base
	var username string
	err = DB.QueryRow("SELECT username FROM users WHERE id = ?", userID).Scan(&username)
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
