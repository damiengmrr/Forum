package handlers

import (
	"database/sql"
	"fmt"
	"forum/database"
	"forum/models"
	"log"
	"net/http"
	"strconv"
	"html/template"

	"github.com/gorilla/sessions"
	"golang.org/x/crypto/bcrypt"
)

var store = sessions.NewCookieStore([]byte("super-secret-key"))

// on stocke la db ici une fois pour toutes (deprecated, on passe par database.GetDatabase())
var DB *sql.DB

func SetDatabase(db *sql.DB) {
	DB = db
}

// ========================= ECHEC HANDLER =========================
func EchecHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "templates/echec.html")
}

// ========================= REGISTER =========================
func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		tmpl, _ := template.ParseFiles("templates/register.html")
		tmpl.Execute(w, nil)
		return
	}

	if r.Method == http.MethodPost {
		username := r.FormValue("username")
		email := r.FormValue("email")
		password := r.FormValue("password")

		if username == "" || email == "" || password == "" {
			log.Println("⚠️ Champs manquants pour l'inscription")
			http.ServeFile(w, r, "templates/ErrorRegister.html")
			return
		}

		hashedPwd, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			log.Println("❌ Erreur hash mot de passe:", err)
			http.ServeFile(w, r, "templates/ErrorRegister.html")
			return
		}

		db := database.GetDatabase()
		if db == nil {
			log.Println("❌ Base de données non initialisée pour register")
			http.Redirect(w, r, "/echec", http.StatusSeeOther)
			return
		}

		_, err = db.Exec("INSERT INTO users (username, email, password) VALUES (?, ?, ?)", username, email, string(hashedPwd))
		if err != nil {
			log.Println("❌ Erreur insertion utilisateur:", err)
			http.ServeFile(w, r, "templates/ErrorRegister.html")
			return
		}

		log.Println("✅ Nouvel utilisateur enregistré:", username)
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
			log.Println("❌ Base de données non initialisée pour login")
			http.Redirect(w, r, "/echec", http.StatusSeeOther)
			return
		}

		var id int
		var username, hashedPassword string

		err := db.QueryRow("SELECT id, username, password FROM users WHERE email = ?", email).Scan(&id, &username, &hashedPassword)
		if err != nil {
			log.Println("❌ Email inconnu ou erreur DB:", err)
			http.ServeFile(w, r, "templates/ErrorLogin.html")
			return
		}

		err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
		if err != nil {
			log.Println("❌ Mot de passe incorrect pour:", email)
			http.ServeFile(w, r, "templates/ErrorLogin.html")
			return
		}

		// création des cookies session
		http.SetCookie(w, &http.Cookie{
			Name:  "username",
			Value: username,
			Path:  "/",
		})

		http.SetCookie(w, &http.Cookie{
			Name:  "session",
			Value: strconv.Itoa(id),
			Path:  "/",
		})

		log.Println("✅ Connexion réussie pour:", username)
		http.Redirect(w, r, "/home", http.StatusSeeOther)
	}
}

// ========================= LOGOUT =========================
func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:   "session",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	})

	http.Redirect(w, r, "/home", http.StatusSeeOther)
}

// ========================= GET CURRENT USER =========================
func GetCurrentUser(r *http.Request) (int, string, error) {
	cookie, err := r.Cookie("session")
	if err != nil {
		return 0, "", fmt.Errorf("pas de session")
	}

	userID := cookie.Value

	db := database.GetDatabase()
	if db == nil {
		return 0, "", fmt.Errorf("base non initialisée")
	}

	var username string
	err = db.QueryRow("SELECT username FROM users WHERE id = ?", userID).Scan(&username)
	if err != nil {
		return 0, "", fmt.Errorf("user introuvable")
	}

	id, err := strconv.Atoi(userID)
	if err != nil {
		return 0, "", fmt.Errorf("id invalide")
	}

	return id, username, nil
}

// ========================= HOME =========================
type PostWithFormattedDate struct {
	models.Post
	FormattedDate string
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("username")
	username := "Invité"
	if err == nil && cookie.Value != "" {
		username = cookie.Value
	}

	rawPosts, err := database.GetAllPosts()
	if err != nil {
		log.Println("❌ Erreur récupération posts home:", err)
		http.Redirect(w, r, "/echec", http.StatusSeeOther)
		return
	}

	var posts []PostWithFormattedDate
	for _, post := range rawPosts {
		posts = append(posts, PostWithFormattedDate{
			Post:          post,
			FormattedDate: post.Date.Format("02 Jan 2006 à 15:04"),
		})
	}

	data := struct {
		Username string
		Posts    []PostWithFormattedDate
		LoggedIn bool
	}{
		Username: username,
		Posts:    posts,
		LoggedIn: username != "Invité",
	}

	tmpl, err := template.ParseFiles("templates/home.html")
	if err != nil {
		log.Println("❌ Erreur template home:", err)
		http.Redirect(w, r, "/echec", http.StatusSeeOther)
		return
	}

	tmpl.Execute(w, data)
}

// ========================= ACCOUNT =========================
func AccountHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session")
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	db := database.GetDatabase()
	if db == nil {
		log.Println("❌ Base non initialisée pour account")
		http.Redirect(w, r, "/echec", http.StatusSeeOther)
		return
	}

	var username, email string
	err = db.QueryRow("SELECT username, email FROM users WHERE id = ?", cookie.Value).Scan(&username, &email)
	if err != nil {
		log.Println("❌ Erreur récupération user account:", err)
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
		log.Println("❌ Erreur template account:", err)
		http.Redirect(w, r, "/echec", http.StatusSeeOther)
		return
	}

	tmpl.Execute(w, data)
}

// ========================= SETTINGS / CONTACT =========================
func SettingsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		http.ServeFile(w, r, "templates/settings.html")
	}
}

func ContactHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		http.ServeFile(w, r, "templates/contact.html")
	}
}
