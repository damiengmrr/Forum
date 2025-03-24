package database

import (
	"database/sql"
	"forum/models"
	"strings"
	"time"
)

var DB *sql.DB

// Stockage en mémoire pour simplifier
var Posts = []models.Post{
	{
		ID:         1,
		Author:     "NonoPrime",
		Title:      "GNIEEEHH",
		Content:    "On vas à l'entrepôt vendredi ?",
		Categories: []string{"Discussion"},
		Date:       time.Now(),
		Status:     "published",
	},
	{
		ID:         2,
		Author:     "TetePrime",
		Title:      "Nouvelles technologies",
		Content:    "Les nouvelles technologies avancent vite, qu'en pensez-vous ?",
		Categories: []string{"Technologie"},
		Date:       time.Now(),
		Status:     "published",
	},
	{
		ID:         3,
		Author:     "Charlie",
		Title:      "Vos jeux préférés",
		Content:    "Quel est votre jeu préféré en ce moment ?",
		Categories: []string{"Jeux vidéo"},
		Date:       time.Now(),
		Status:     "draft",
	},
}

// Fonction pour ajouter un post
func AddPost(author, title, content, categories string) {
	cats := strings.Split(categories, ",")
	for i := range cats {
		cats[i] = strings.TrimSpace(cats[i])
	}

	newPost := models.Post{
		ID:         len(Posts) + 1,
		Author:     author,
		Title:      title,
		Content:    content,
		Categories: cats,
		Date:       time.Now(),
		Status:     "draft", // Par défaut en brouillon
	}
	Posts = append(Posts, newPost)
}

// Fonction pour changer le statut d'un post
func UpdatePostStatus(id int, status string) {
	for i, p := range Posts {
		if p.ID == id {
			Posts[i].Status = status
			break
		}
	}
}

// Fonction pour supprimer un post
func DeletePost(id int) {
	for i, p := range Posts {
		if p.ID == id {
			Posts = append(Posts[:i], Posts[i+1:]...)
			break
		}
	}
}
