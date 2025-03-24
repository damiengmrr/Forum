package handlers

import (
	"fmt"
	"forum/models"
	"html/template"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/google/uuid"
)

// Slice globale des posts utilisée pour l'affichage sur la page d'accueil
var posts = []models.Post{
	{
		ID:         1,
		Author:     "John",
		Title:      "Test Post",
		Content:    "Ceci est un post de test",
		Date:       time.Now(),
		Likes:      5,
		Dislikes:   2,
		Categories: []string{"discussion"},
		ImagePath:  "",
	},
	{
		ID:         2,
		Author:     "Alice",
		Title:      "Technologie",
		Content:    "Post sur la technologie",
		Date:       time.Now(),
		Likes:      8,
		Dislikes:   1,
		Categories: []string{"technology"},
		ImagePath:  "",
	},
}

// Handler pour afficher l'heure formatée
func TimeHandlers(w http.ResponseWriter, r *http.Request) {
	currentTime := time.Now()
	formattedTime := currentTime.Format("02 Jan 2006 à 15:04")
	fmt.Fprintf(w, "Date : %s", formattedTime)
}

// CreatePostHandler gère l'affichage du formulaire (GET) et la création d'un nouveau post (POST)
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

	if r.Method == "POST" {
		// Parse le formulaire multipart avec une limite de 20 MB
		if err := r.ParseMultipartForm(20 << 20); err != nil {
			http.Error(w, "Erreur lors du parsing du formulaire", http.StatusInternalServerError)
			return
		}

		// Récupération des champs texte
		title := r.FormValue("title")
		author := r.FormValue("author")
		content := r.FormValue("content")
		// Récupération de la sélection multiple pour les catégories
		categories := r.Form["categories"]

		// Récupération du fichier image (si fourni)
		file, handler, err := r.FormFile("image")
		var imagePath string
		if err == nil && handler.Size > 0 {
			defer file.Close()

			// Génère un nom unique pour l'image grâce à UUID et conserve l'extension originale
			ext := filepath.Ext(handler.Filename)
			newFileName := uuid.New().String() + ext

			// Crée le fichier dans le dossier static/uploads (assure-toi que le dossier existe)
			out, err := os.Create("./static/uploads/" + newFileName)
			if err != nil {
				fmt.Println("Erreur lors de la création du fichier :", err)
				http.Error(w, "Erreur lors de la création du fichier sur le serveur", http.StatusInternalServerError)
				return
			}
			defer out.Close()

			// Copie le contenu de l'image uploadée dans le fichier créé
			if _, err = io.Copy(out, file); err != nil {
				http.Error(w, "Erreur lors de la copie du fichier", http.StatusInternalServerError)
				return
			}

			// Stocke le chemin relatif de l'image
			imagePath = "/static/uploads/" + newFileName
		}

		// Création du nouveau post
		newPost := models.Post{
			ID:         len(posts) + 1, // ou utiliser un mécanisme d'ID plus robuste
			Title:      title,
			Author:     author,
			Content:    content,
			Date:       time.Now(),
			Likes:      0,
			Dislikes:   0,
			Categories: categories,
			ImagePath:  imagePath,
		}

		fmt.Printf("Nouveau post créé : %+v\n", newPost)

		// Ajoute le nouveau post à la slice globale
		posts = append(posts, newPost)

		// Redirige vers la page d'accueil après création
		http.Redirect(w, r, "/home", http.StatusSeeOther)
		return
	}

	// Par défaut, renvoie le template de création de post (GET)
	tmpl := template.Must(template.ParseFiles("templates/create-post.html"))
	tmpl.Execute(w, nil)
}

// PostHandler affiche un post spécifique
func PostHandler(w http.ResponseWriter, r *http.Request) {
	// Récupère l'ID du post depuis l'URL (ex: /post/1)
	idStr := r.URL.Path[len("/post/"):]
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 || id > len(posts) {
		http.Error(w, "Post non trouvé", http.StatusNotFound)
		return
	}

	// Récupère le post correspondant (les IDs commencent à 1)
	post := posts[id-1]

	// Charge et exécute le template du post
	tmpl, err := template.ParseFiles("templates/post.html")
	if err != nil {
		http.Error(w, "Erreur serveur", http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, post)
}

// LikeHandler gere le like d'un post avec limitation via cookies
func LikeHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID invalide", http.StatusBadRequest)
		return
	}

	cookieName := fmt.Sprintf("vote_post_%d", id)
	voteCookie, err := r.Cookie(cookieName)
	vote := ""
	if err == nil {
		vote = voteCookie.Value
	}

	for i := range posts {
		if posts[i].ID == id {
			if vote == "like" {
				// déjà liké → rien
			} else if vote == "dislike" {
				posts[i].Dislikes--
				posts[i].Likes++
				http.SetCookie(w, &http.Cookie{Name: cookieName, Value: "like", Path: "/"})
			} else {
				posts[i].Likes++
				http.SetCookie(w, &http.Cookie{Name: cookieName, Value: "like", Path: "/"})
			}
			break
		}
	}

	http.Redirect(w, r, "/home", http.StatusSeeOther)
}

// DislikeHandler gere le dislike d'un post avec limitation via cookies
func DislikeHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID invalide", http.StatusBadRequest)
		return
	}

	cookieName := fmt.Sprintf("vote_post_%d", id)
	voteCookie, err := r.Cookie(cookieName)
	vote := ""
	if err == nil {
		vote = voteCookie.Value
	}

	for i := range posts {
		if posts[i].ID == id {
			if vote == "dislike" {
				// déjà disliké → rien
			} else if vote == "like" {
				posts[i].Likes--
				posts[i].Dislikes++
				http.SetCookie(w, &http.Cookie{Name: cookieName, Value: "dislike", Path: "/"})
			} else {
				posts[i].Dislikes++
				http.SetCookie(w, &http.Cookie{Name: cookieName, Value: "dislike", Path: "/"})
			}
			break
		}
	}

	http.Redirect(w, r, "/home", http.StatusSeeOther)
}