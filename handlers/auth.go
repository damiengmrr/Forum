package handlers
import (
	"database/sql"
	"fmt"
	"net/http"
	
	"forum/database"
	"golang.org/x/crypto/bcrypt"
	"github.com/gofrs/uuid"
	)
	
	// RegisterHandler gère l'inscription des utilisateurs
	func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
	username := r.FormValue("username")
	email := r.FormValue("email")
	password := r.FormValue("password")
	}
}
	
	// Vérifier si l'email existe déjà
	var exists string
	err := database.DB.QueryRow("SELECT email FROM users WHERE email = ?", email).Scan(&exists)
	if err == nil {
	http.Error(w, "Email déjà utilisé", http.StatusBadRequest)
	return
	}
	
	// Hasher le mot de passe
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
	http.Error(w, "Erreur lors du hash du mot de passe", http.StatusInternalServerError)
	return
	}
	
	// Générer un UUID pour l'utilisateur
	id, _ := uuid.NewV4()
	
	// Insérer l'utilisateur en base
	_, err = database.DB.Exec("INSERT INTO users (id, username, email, password) VALUES (?, ?, ?, ?)",
	id.String(), username, email, string(hashedPassword))
	if err != nil {
	http.Error(w, "Erreur lors de l'inscription", http.StatusInternalServerError)
	...(91lignes restantes)
	}