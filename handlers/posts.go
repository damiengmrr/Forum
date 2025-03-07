package handlers

import (
	"net/http"
	"time"

	"forum/models"
	"forum/database"
)

// Gestionnaire pour créer un nouveau post
func CreatePostHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
		return
	}

	author := r.FormValue("author")
	content := r.FormValue("content")
	category := r.FormValue("category")

	if author == "" || content == "" || category == "" {
		http.Error(w, "Tous les champs sont requis", http.StatusBadRequest)
		return
	}

	// Ajouter un nouveau post dans database.Posts
	newPost := models.Post{
		ID:       len(database.Posts) + 1,
		Author:   author,
		Content:  content,
		Category: category,
		Date:     time.Now(),
	}
	database.Posts = append(database.Posts, newPost)

	// Rediriger vers la page d'accueil pour voir le post ajouté
	http.Redirect(w, r, "/home", http.StatusSeeOther)
}