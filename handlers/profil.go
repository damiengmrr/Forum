package handlers

import (
	"forum/database"
	"html/template"
	"log"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

func EditProfileHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		log.Println("📝 Accès à la page de modification du profil")
		tmpl, _ := template.ParseFiles("templates/edit-profile.html")
		tmpl.Execute(w, nil)
		return
	}

	if r.Method == "POST" {
		log.Println("📩 Tentative de modification du pseudo")

		cookie, err := r.Cookie("session")
		if err != nil {
			log.Println("❌ Pas de cookie de session, redirection vers login")
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		newUsername := r.FormValue("username")
		if newUsername == "" {
			log.Println("⚠️ Aucun nouveau pseudo fourni")
			http.Redirect(w, r, "/edit-profile", http.StatusSeeOther)
			return
		}

		db := database.GetDatabase()
		_, err = db.Exec("UPDATE users SET username = ? WHERE id = ?", newUsername, cookie.Value)
		if err != nil {
			log.Println("❌ Erreur lors de la mise à jour du pseudo :", err)
			http.Redirect(w, r, "/echec", http.StatusSeeOther)
			return
		}

		// ✅ Mise à jour du cookie
		http.SetCookie(w, &http.Cookie{
			Name:  "username",
			Value: newUsername,
			Path:  "/",
		})

		log.Println("✅ Pseudo mis à jour avec succès :", newUsername)

		http.Redirect(w, r, "/account", http.StatusSeeOther)
	}
}

func ChangePasswordHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		log.Println("📝 Accès à la page de changement de mot de passe")
		tmpl, _ := template.ParseFiles("templates/change-password.html")
		tmpl.Execute(w, nil)
		return
	}

	if r.Method == "POST" {
		log.Println("📩 Tentative de changement de mot de passe")

		cookie, err := r.Cookie("session")
		if err != nil {
			log.Println("❌ Pas de cookie de session, redirection vers login")
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		currentPassword := r.FormValue("currentPassword")
		newPassword := r.FormValue("newPassword")

		db := database.GetDatabase()

		var hashedPassword string
		err = db.QueryRow("SELECT password FROM users WHERE id = ?", cookie.Value).Scan(&hashedPassword)
		if err != nil {
			log.Println("❌ Erreur récupération mot de passe depuis la base :", err)
			http.Redirect(w, r, "/echec", http.StatusSeeOther)
			return
		}

		err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(currentPassword))
		if err != nil {
			log.Println("⚠️ Mot de passe actuel incorrect")
			http.Redirect(w, r, "/echec", http.StatusSeeOther)
			return
		}

		newHashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
		if err != nil {
			log.Println("❌ Erreur lors du hash du nouveau mot de passe :", err)
			http.Redirect(w, r, "/echec", http.StatusSeeOther)
			return
		}

		_, err = db.Exec("UPDATE users SET password = ? WHERE id = ?", string(newHashedPassword), cookie.Value)
		if err != nil {
			log.Println("❌ Erreur lors de la mise à jour du mot de passe dans la base :", err)
			http.Redirect(w, r, "/echec", http.StatusSeeOther)
			return
		}

		log.Println("✅ Mot de passe changé avec succès pour l'utilisateur ID :", cookie.Value)
		http.Redirect(w, r, "/account", http.StatusSeeOther)
	}
}
