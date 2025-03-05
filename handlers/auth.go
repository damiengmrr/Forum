package handlers

import (
	"database/sql"
	"fmt"
	"net/http"

	"forum/database"

	//"github.com/gofrs/uuid"
	"golang.org/x/crypto/bcrypt"
)

// RegisterHandler gère l'inscription des utilisateurs
func RegisterHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method == "POST" {
        email := r.FormValue("email")
        password := r.FormValue("password")

        // Vérifier si l'email existe déjà
        var storedEmail string
        err := database.DB.QueryRow("SELECT email FROM users WHERE email = ?", email).Scan(&storedEmail)
        if err == nil {
            // Si l'email existe déjà, redirige vers la page d'erreur
            http.ServeFile(w, r, "templates/ErrorRegister.html")
            return
        } else if err != sql.ErrNoRows {
            // Autre erreur de base de données
            http.Error(w, "Erreur serveur", http.StatusInternalServerError)
            return
        }

        // Hacher le mot de passe avant de l'enregistrer
        hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
        if err != nil {
            http.ServeFile(w, r, "templates/ErrorRegister.html")
            return
        }

        // Enregistrer l'utilisateur dans la base de données
        _, err = database.DB.Exec("INSERT INTO users (email, password) VALUES (?, ?)", email, string(hashedPassword))
        if err != nil {
            http.ServeFile(w, r, "templates/ErrorRegister.html")
            return
        }

        fmt.Fprintln(w, "Inscription réussie !")
        return
    } else {
        http.ServeFile(w, r, "templates/register.html")
    }
}

// LoginHandler gère la connexion des utilisateurs
func LoginHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method == "POST" {
        email := r.FormValue("email")
        password := r.FormValue("password")

        // Vérifier si l'email existe
        var storedHash string
        err := database.DB.QueryRow("SELECT password FROM users WHERE email = ?", email).Scan(&storedHash)
        if err == sql.ErrNoRows {
            // Si l'email n'existe pas, redirige vers la page d'erreur
            http.ServeFile(w, r, "templates/ErrorLogin.html")
            return
        } else if err != nil {
            // Si une autre erreur se produit
            http.Error(w, "Erreur serveur", http.StatusInternalServerError)
            return
        }

        // Vérifier si le mot de passe est correct
        err = bcrypt.CompareHashAndPassword([]byte(storedHash), []byte(password))
        if err != nil {
            // Si le mot de passe est incorrect, redirige vers la page d'erreur
            http.ServeFile(w, r, "templates/ErrorLogin.html")
            return
        }

        // Création d'un cookie de session
        cookie := http.Cookie{
            Name:  "session",
            Value: email, // Simple pour le moment, à améliorer
            Path:  "/",
        }
        http.SetCookie(w, &cookie)

        // Redirection vers la page d'accueil après connexion réussie
        http.Redirect(w, r, "/home", http.StatusFound)
        return
    } else {
        // Si la méthode n'est pas POST, on affiche le formulaire de connexion
        http.ServeFile(w, r, "templates/login.html")
    }
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		http.ServeFile(w, r, "templates/home.html")
	}
}

func AccountHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
        http.ServeFile(w, r, "templates/account.html")
    }
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method == "GET" {
        http.ServeFile(w, r, "templates/logout.html")
}
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