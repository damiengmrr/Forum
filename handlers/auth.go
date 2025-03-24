package handlers

import (
	"database/sql"
	"fmt"
	"forum/database"
	"net/http"

	//"github.com/gofrs/uuid"
	"html/template"

	"golang.org/x/crypto/bcrypt"
)

func EchecHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "templates/echec.html")
}

// RegisterHandler gère l'inscription des utilisateurs
func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		// Récupérer les informations du formulaire
		username := r.FormValue("username")
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
			http.ServeFile(w, r, "templates/echec.html")
			println(err)
			return
		}

		// Hacher le mot de passe avant de l'enregistrer
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			http.ServeFile(w, r, "templates/ErrorRegister.html")
			return
		}

		// Enregistrer l'utilisateur dans la base de données sans spécifier l'ID (ID généré automatiquement par MySQL)
		_, err = database.DB.Exec("INSERT INTO users (username, email, password) VALUES (?, ?, ?)", username, email, string(hashedPassword))
		if err != nil {
			http.ServeFile(w, r, "templates/ErrorRegister.html")
			return
		}

		// Afficher un message de succès
		fmt.Fprintln(w, "Inscription réussie !")
		return
	} else {
		// Si ce n'est pas une requête POST, afficher le formulaire d'inscription
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
			http.ServeFile(w, r, "templates/echec.html")
			println(err)
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

// Afficher la page d'accueil avec tous les posts
func HomeHandler(w http.ResponseWriter, r *http.Request) {
    tmpl, err := template.ParseFiles("templates/home.html")
    if err != nil {
        http.ServeFile(w, r, "templates/echec.html")
			println(err)
		return
    }

    // Envoyer tous les posts dans le template
    //tmpl.Execute(w, database.Posts)
	tmpl.Execute(w, posts) // Ici, on passe la slice globale "posts"
	
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
