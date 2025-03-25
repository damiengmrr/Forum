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
	"golang.org/x/exp/rand"
)

// Slice globale des posts utilisée pour l'affichage sur la page d'accueil
var posts = []models.Post{
	{
		ID:         1,
		Author:     "NonoPrime",
		Title:      "GNIEEEHHH",
		Content:    "On sort à l'entrepôt Vendredi ?",
		Date:       time.Now(),
		Likes:      10,
		Dislikes:   2,
		Categories: []string{"discussion"},
		ImagePath:  "",
		Status:     "published",
		Comments: []models.Comment{
			{
				ID:       2,
				Author:   "TetePrime",
				Avatar:   "/static/image/pfp1.png",
				Content:  "Logique non ?",
				Likes:    3,
				Dislikes: 1,
				Response: &models.Comment{
					ID:       3,
					Author:   "VavaPrime",
					Avatar:   "/static/image/pfp2.png",
					Content:  "Voilaaa théo",
					Likes:    5,
					Dislikes: 0,
				},
			},
			{
				ID:       4,
				Author:   "DadaPrime",
				Avatar:   "/static/image/pfp3.png",
				Content:  "Normal, jeudi MILK ?",
				Likes:    12,
				Dislikes: 8,
				Response: nil,
			},
		},
	},
	{
		ID:         5,
		Author:     "les mentors à leurs prime",
		Title:      "On vous voit hein !",
		Content:    "Vous avez intérêt à rendre vos projets à temps !",
		Date:       time.Now(),
		Likes:      10,
		Dislikes:   2,
		Categories: []string{"discussion"},
		ImagePath:  "",
		Status:     "published",
		Comments: []models.Comment{
			{
				ID:       6,
				Author:   "TetePrime",
				Avatar:   "/static/image/pfp1.png",
				Content:  "tkt tkt",
				Likes:    3,
				Dislikes: 1,
				Response: &models.Comment{
					ID:       7,
					Author:   "DadaPrime",
					Avatar:   "/static/image/pfp2.png",
					Content:  "Daronned sur forum carrement",
					Likes:    5,
					Dislikes: 0,
				},
			},
		},
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
			http.ServeFile(w, r, "templates/echec.html")
			fmt.Print(err)
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
	if err != nil || id <= 0 {
		http.ServeFile(w, r, "templates/echec.html")
		fmt.Print(err)
	}

	// Récupère le post correspondant (les IDs commencent à 1)
	var post models.Post
	found := false
	for _, p := range posts {
		if p.ID == id {
			post = p
			found = true
			break
		}
	}

	if !found {
		http.ServeFile(w, r, "templates/echec.html")
		fmt.Print(err)
	}

	// Charge et exécute le template du post
	tmpl, err := template.ParseFiles("templates/post.html")
	if err != nil {
		http.ServeFile(w, r, "templates/echec.html")
		fmt.Print(err)
	}
	tmpl.Execute(w, post)
}

// LikeHandler gere le like d'un post avec limitation via cookies
func LikeHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.ServeFile(w, r, "templates/echec.html")
		fmt.Print(err)
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
		http.ServeFile(w, r, "templates/echec.html")
		fmt.Print(err)
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
func CommentReplyHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.ServeFile(w, r, "templates/echec.html")
	}

	idStr := r.FormValue("id")
	content := r.FormValue("content")
	if idStr == "" || content == "" {
		http.ServeFile(w, r, "templates/echec.html")
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.ServeFile(w, r, "templates/echec.html")
		fmt.Print(err)
	}

	for pi := range posts {
		for ci := range posts[pi].Comments {
			if posts[pi].Comments[ci].ID == id && posts[pi].Comments[ci].Response == nil {
				posts[pi].Comments[ci].Response = &models.Comment{
					ID:       rand.Intn(100000), // ID random simple
					Author:   "Toi",
					Avatar:   "/static/image/default.png",
					Content:  content,
					Likes:    0,
					Dislikes: 0,
				}
				break
			}
		}
	}

	http.Redirect(w, r, r.Referer(), http.StatusSeeOther)
}

// ici on gere le like d'un commentaire
func CommentLikeHandler(w http.ResponseWriter, r *http.Request) {
	// on recupere l'id du commentaire dans l'URL
	idStr := r.URL.Query().Get("id")
	id, _ := strconv.Atoi(idStr)

	// on utilise un cookie pour savoir si l'utilisateur a deja vote
	cookieName := fmt.Sprintf("comment_vote_%d", id)
	vote := ""

	// si le cookie existe deja, on le lit
	if c, err := r.Cookie(cookieName); err == nil {
		vote = c.Value
	}

	// on cherche dans tous les posts et leurs commentaires
	for pi := range posts {
		for ci := range posts[pi].Comments {
			// si on trouve le commentaire
			if posts[pi].Comments[ci].ID == id {
				if vote == "like" {
					// il a deja like → on fait rien
				} else if vote == "dislike" {
					// il avait dislike → on inverse
					posts[pi].Comments[ci].Dislikes--
					posts[pi].Comments[ci].Likes++
				} else {
					// il avait rien fait → on like
					posts[pi].Comments[ci].Likes++
				}
				// on met à jour le cookie avec "like"
				http.SetCookie(w, &http.Cookie{Name: cookieName, Value: "like", Path: "/"})
			}

			// si on like une réponse au commentaire
			if posts[pi].Comments[ci].Response != nil && posts[pi].Comments[ci].Response.ID == id {
				if vote == "like" {
				} else if vote == "dislike" {
					posts[pi].Comments[ci].Response.Dislikes--
					posts[pi].Comments[ci].Response.Likes++
				} else {
					posts[pi].Comments[ci].Response.Likes++
				}
				http.SetCookie(w, &http.Cookie{Name: cookieName, Value: "like", Path: "/"})
			}
		}
	}

	http.Redirect(w, r, r.Referer(), http.StatusSeeOther)
}

// ici on gere le dislike d'un commentaire
func CommentDislikeHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, _ := strconv.Atoi(idStr)

	cookieName := fmt.Sprintf("comment_vote_%d", id)
	vote := ""

	if c, err := r.Cookie(cookieName); err == nil {
		vote = c.Value
	}

	for pi := range posts {
		for ci := range posts[pi].Comments {
			if posts[pi].Comments[ci].ID == id {
				if vote == "dislike" {
					// deja dislike → on fait rien
				} else if vote == "like" {
					posts[pi].Comments[ci].Likes--
					posts[pi].Comments[ci].Dislikes++
				} else {
					posts[pi].Comments[ci].Dislikes++
				}
				http.SetCookie(w, &http.Cookie{Name: cookieName, Value: "dislike", Path: "/"})
			}

			if posts[pi].Comments[ci].Response != nil && posts[pi].Comments[ci].Response.ID == id {
				if vote == "dislike" {
				} else if vote == "like" {
					posts[pi].Comments[ci].Response.Likes--
					posts[pi].Comments[ci].Response.Dislikes++
				} else {
					posts[pi].Comments[ci].Response.Dislikes++
				}
				http.SetCookie(w, &http.Cookie{Name: cookieName, Value: "dislike", Path: "/"})
			}
		}
	}

	http.Redirect(w, r, r.Referer(), http.StatusSeeOther)
}
