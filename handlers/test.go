package handlers

import (
	"fmt"
	"net/http"
	"time"

	"forum/database"
)

// TestSessionHandler affiche toutes les infos utiles pour debug
func TestSessionHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")

	//fmt.Fprintln(w, "=== DEBUG : INFOS DE SESSION EN COURS ===\n")

	// 🔍 Cookies
	fmt.Fprintln(w, "🍪 Cookies reçus :")
	for _, c := range r.Cookies() {
		fmt.Fprintf(w, "- %s = %s\n", c.Name, c.Value)
	}
	fmt.Fprintln(w)

	// 👤 Utilisateur connecté
	userID, username, err := GetCurrentUser(r)
	if err != nil {
		fmt.Fprintln(w, "⚠️ Utilisateur non connecté :", err)
	} else {
		fmt.Fprintln(w, "👤 Utilisateur connecté :")
		fmt.Fprintf(w, "- ID       : %d\n", userID)
		fmt.Fprintf(w, "- Username : %s\n", username)
	}
	fmt.Fprintln(w)

	// 🌍 Adresse IP
	fmt.Fprintln(w, "🌍 Adresse IP (RemoteAddr) :", r.RemoteAddr)
	fmt.Fprintln(w)

	// 🔡 Headers HTTP
	fmt.Fprintln(w, "📦 Headers HTTP :")
	for name, values := range r.Header {
		for _, value := range values {
			fmt.Fprintf(w, "- %s: %s\n", name, value)
		}
	}
	fmt.Fprintln(w)

	// 📄 Méthode et URL demandée
	fmt.Fprintln(w, "🧭 Requête actuelle :")
	fmt.Fprintf(w, "- Méthode : %s\n", r.Method)
	fmt.Fprintf(w, "- URL     : %s\n", r.URL.Path)
	fmt.Fprintf(w, "- Query   : %s\n", r.URL.RawQuery)
	fmt.Fprintln(w)

	// 📱 User-Agent
	fmt.Fprintln(w, "📱 User-Agent :")
	fmt.Fprintln(w, r.UserAgent())
	fmt.Fprintln(w)

	// 🕒 Timestamp
	fmt.Fprintln(w, "🕒 Date actuelle :")
	fmt.Fprintln(w, time.Now().Format("02 Jan 2006 15:04:05"))
	fmt.Fprintln(w)

	// 🗄️ Infos base de données
	fmt.Fprintln(w, "🗄️ Base de données :")
	db := database.GetDatabase()
	if db == nil {
		fmt.Fprintln(w, "❌ Connexion BDD : NON ÉTABLIE")
	} else {
		err = db.Ping()
		if err != nil {
			fmt.Fprintln(w, "⚠️ Connexion BDD : Erreur de ping :", err)
		} else {
			fmt.Fprintln(w, "✅ Connexion BDD : OK")

			// Compter les utilisateurs
			var userCount int
			err = db.QueryRow("SELECT COUNT(*) FROM users").Scan(&userCount)
			if err != nil {
				fmt.Fprintln(w, "⚠️ Erreur comptage utilisateurs :", err)
			} else {
				fmt.Fprintf(w, "- Nombre d'utilisateurs : %d\n", userCount)
			}

			// Compter les posts
			var postCount int
			err = db.QueryRow("SELECT COUNT(*) FROM posts").Scan(&postCount)
			if err != nil {
				fmt.Fprintln(w, "⚠️ Erreur comptage posts :", err)
			} else {
				fmt.Fprintf(w, "- Nombre de posts : %d\n", postCount)
			}
		}
	}
	fmt.Fprintln(w)

	// 📝 Logs ou erreurs récentes (placeholder pour l'instant)
	fmt.Fprintln(w, "📝 Logs / Erreurs récentes :")
	fmt.Fprintln(w, "- (Pas de système de logs actif actuellement)")
	fmt.Fprintln(w)

	fmt.Fprintln(w, "✅ Fin du test.")
}
