package handlers

import (
	"html/template"
	"net/http"
	"strconv"
	"time"

	"forum/database"
	"forum/models"
)
var posts = []models.Post{
    {ID: 1, Author: "John", Content: "Ceci est un post de test", Category: "Discussion", Date: time.Now(), Likes: 5, Dislikes: 2},
    {ID: 2, Author: "Alice", Content: "Post sur la technologie", Category: "Technologie", Date: time.Now(), Likes: 8, Dislikes: 1},
}
// Gestionnaire pour créer un nouveau post
func CreatePostHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		tmpl, err := template.ParseFiles("templates/create-post.html")
		if err != nil {
			http.Error(w, "Erreur lors du chargement de la page", http.StatusInternalServerError)
			return
		}
		tmpl.Execute(w, nil)
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
func PostHandler(w http.ResponseWriter, r *http.Request) {
	// Récupérer l'ID du post depuis l'URL
	idStr := r.URL.Path[len("/post/"):]
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 || id > len(posts) {
		http.Error(w, "Post non trouvé", http.StatusNotFound)
		return
	}

	// Trouver le post correspondant
	post := posts[id-1]

	// Affichage de la page du post
	tmpl, err := template.ParseFiles("templates/post.html")
	if err != nil {
		http.Error(w, "Erreur serveur", http.StatusInternalServerError)
		return
	}

	tmpl.Execute(w, post)
}

/* BDD avec
package handlers

import (
    "html/template"
    "net/http"
    "forum/models"
    "forum/database"
)

func CreatePostHandler(w http.ResponseWriter, r *http.Request) {
    // Si la méthode est GET, afficher le formulaire
    if r.Method == "GET" {
        tmpl, err := template.ParseFiles("templates/create-post.html")
        if err != nil {
            http.Error(w, "Erreur lors du chargement de la page", http.StatusInternalServerError)
            return
        }
        tmpl.Execute(w, nil)
        return
    }

    // Si la méthode est POST, traiter les données du formulaire
    if r.Method == "POST" {
        // Récupérer les données du formulaire
        author := r.FormValue("author")
        content := r.FormValue("content")
        category := r.FormValue("category")

        // Créer un nouveau post
        post := models.Post{
            Author:   author,
            Content:  content,
            Category: category,
            Date:     time.Now(),
        }

        // Sauvegarder le post dans la base de données
        err := database.CreatePost(post)
        if err != nil {
            http.Error(w, "Erreur lors de la création du post", http.StatusInternalServerError)
            return
        }

        // Rediriger vers la page d'accueil après la création du post
        http.Redirect(w, r, "/home", http.StatusSeeOther)
    }
}
*/
