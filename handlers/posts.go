package handlers

import (
	"database/sql"
	"fmt"
	"forum/database"
	"forum/models"
	"html/template"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
	"time"
)

// pour afficher un post + ses commentaires
func PostHandler(w http.ResponseWriter, r *http.Request) {
	// ✅ récupère l'utilisateur connecté
	_, username, _ := GetCurrentUser(r)

	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Println("❌ id invalide :", idStr)
		http.Redirect(w, r, "/echec", http.StatusSeeOther)
		return
	}

	post, err := database.GetPostByID(id)
	if err != nil {
		log.Println("❌ post introuvable :", err)
		http.Redirect(w, r, "/echec", http.StatusSeeOther)
		return
	}

	// ✅ récupère la photo de profil de l'auteur du post
	var profilePicture sql.NullString
	err = database.GetDatabase().QueryRow("SELECT profile_picture FROM users WHERE username = ?", post.Author).Scan(&profilePicture)
	if err != nil {
		log.Println("❌ erreur récupération photo profil auteur :", err)
	}
	pic := "default.jpg"
	if profilePicture.Valid && profilePicture.String != "" {
		pic = profilePicture.String
	}

	comments, err := database.GetCommentsByPostID(id)
	if err != nil {
		log.Println("❌ erreur récup commentaires :", err)
		http.Redirect(w, r, "/echec", http.StatusSeeOther)
		return
	}

	// 🟢 injecter les commentaires dans le post
	post.Comments = comments

	// ✅ structure enrichie avec la photo de profil
	data := struct {
		Post           models.Post
		FormattedDate  string
		Comments       []models.Comment
		IsAuthor       bool
		ProfilePicture string
	}{
		Post:           post,
		FormattedDate:  post.Date.Format("02 Jan 2006 à 15:04"),
		Comments:       comments,
		IsAuthor:       post.Author == username,
		ProfilePicture: pic,
	}

	tmpl, err := template.ParseFiles("templates/post.html")
	if err != nil {
		log.Println("❌ erreur template :", err)
		http.Redirect(w, r, "/echec", http.StatusSeeOther)
		return
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		log.Println("❌ erreur execute template :", err)
	}
}

// pour filtrer les posts par catégorie
func CategoriesHandler(w http.ResponseWriter, r *http.Request) {
	cat := r.URL.Query().Get("name")

	var posts []models.Post
	var err error
	var formatted []PostWithFormattedDate

	if cat != "" {
		// si on a cliqué sur une catégorie, on filtre
		posts, err = database.GetPostsByCategory(cat)
		if err != nil {
			http.Redirect(w, r, "/echec", http.StatusSeeOther)
			return
		}

		for _, p := range posts {
			formatted = append(formatted, PostWithFormattedDate{
				Post:          p,
				FormattedDate: p.Date.Format("02 Jan 2006 à 15:04"),
			})
		}
	}

	data := struct {
		Categories       []string
		SelectedCategory string
		FilteredPosts    []PostWithFormattedDate
	}{
		Categories:       []string{"Discussion", "Technologie", "Jeux Vidéo", "Littérature"},
		SelectedCategory: cat,
		FilteredPosts:    formatted,
	}

	tmpl, err := template.ParseFiles("templates/categories.html")
	if err != nil {
		http.Redirect(w, r, "/echec", http.StatusSeeOther)
		return
	}

	tmpl.Execute(w, data)
}

// pour créer un nouveau post
func CreatePostHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		tmpl, _ := template.ParseFiles("templates/create-post.html")
		tmpl.Execute(w, nil)
		return
	}

	// on récupère l'utilisateur connecté
	userID, username, err := GetCurrentUser(r)
	if err != nil || userID == 0 {
		log.Println("❌ utilisateur non connecté :", err)
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// méthode POST : traitement du formulaire
	r.ParseMultipartForm(10 << 20)

	title := r.FormValue("title")
	content := r.FormValue("content")
	categories := r.Form["categories"]
	date := time.Now()

	imagePath := "" // par défaut
	file, handler, err := r.FormFile("image")
	if err == nil {
		defer file.Close()
		imagePath = "/static/uploads/" + handler.Filename
		dst, err := saveFile(file, handler)
		if err != nil {
			log.Println("Erreur upload image:", err)
		}
		log.Println("Image enregistrée dans :", dst)
	}
	log.Println("👤 Auteur du post :", username)
	newPost := models.Post{
		Author:     username,
		Title:      title,
		Content:    content,
		Categories: categories,
		Date:       date,
		ImagePath:  imagePath,
		Likes:      0,
		Dislikes:   0,
	}

	err = database.InsertPost(newPost)
	if err != nil {
		log.Println("❌ erreur insert post :", err)
		http.Redirect(w, r, "/echec", http.StatusSeeOther)
		return
	}

	log.Println("🟢 Nouveau post créé par", username, ":", title)
	http.Redirect(w, r, "/home", http.StatusSeeOther)
}

// pour enregistrer le fichier image dans /static/uploads/
func saveFile(file multipart.File, handler *multipart.FileHeader) (string, error) {
	dstPath := "static/uploads/" + handler.Filename
	dst, err := os.Create(dstPath)
	if err != nil {
		return "", err
	}
	defer dst.Close()

	_, err = io.Copy(dst, file)
	return dstPath, err
}
func LikeHandler(w http.ResponseWriter, r *http.Request) {
	userID := GetCurrentUserID(r)
	postIDStr := r.URL.Query().Get("id")
	postID, err := strconv.Atoi(postIDStr)
	if err != nil || userID == 0 {
		http.Redirect(w, r, "/echec", http.StatusSeeOther)
		fmt.Print(err)
		fmt.Print(userID)
		return
	}

	err = database.TogglePostVote(userID, postID, "like")
	if err != nil {
		log.Println("❌ erreur toggle like :", err)
		http.Redirect(w, r, "/echec", http.StatusSeeOther)
		return
	}

	http.Redirect(w, r, r.Referer(), http.StatusSeeOther)
}

func DislikeHandler(w http.ResponseWriter, r *http.Request) {
	userID := GetCurrentUserID(r)
	postIDStr := r.URL.Query().Get("id")
	postID, err := strconv.Atoi(postIDStr)
	if err != nil || userID == 0 {
		http.Redirect(w, r, "/echec", http.StatusSeeOther)
		fmt.Print(err)
		return
	}

	err = database.TogglePostVote(userID, postID, "dislike")
	if err != nil {
		log.Println("❌ erreur toggle dislike :", err)
		http.Redirect(w, r, "/echec", http.StatusSeeOther)
		return
	}

	http.Redirect(w, r, r.Referer(), http.StatusSeeOther)
}

func CommentReplyHandler(w http.ResponseWriter, r *http.Request) {
	userID, username, err := GetCurrentUser(r)
	if err != nil || userID == 0 {
		log.Println("❌ utilisateur non connecté")
		http.Redirect(w, r, "/echec", http.StatusSeeOther)
		return
	}

	// ID du post principal
	postIDStr := r.FormValue("post_id")
	postID, err := strconv.Atoi(postIDStr)
	if err != nil {
		log.Println("❌ id du post invalide")
		http.Redirect(w, r, "/echec", http.StatusSeeOther)
		return
	}

	// ID du commentaire auquel on répond (ou vide/null si c’est une réponse au post)
	responseToStr := r.FormValue("response_to")
	var responseTo sql.NullInt64
	if responseToStr != "" {
		responseID, err := strconv.Atoi(responseToStr)
		if err == nil {
			responseTo = sql.NullInt64{Int64: int64(responseID), Valid: true}
		}
	}

	content := r.FormValue("content")
	if content == "" {
		log.Println("❌ contenu vide")
		http.Redirect(w, r, "/echec", http.StatusSeeOther)
		return
	}

	err = database.InsertComment(postID, username, content, responseTo)
	if err != nil {
		log.Println("❌ erreur insertion commentaire :", err)
		http.Redirect(w, r, "/echec", http.StatusSeeOther)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/post?id=%d", postID), http.StatusSeeOther)
}

func CommentLikeHandler(w http.ResponseWriter, r *http.Request) {
	userID := GetCurrentUserID(r)
	commentIDStr := r.URL.Query().Get("id")
	commentID, err := strconv.Atoi(commentIDStr)
	if err != nil || userID == 0 {
		http.Redirect(w, r, "/echec", http.StatusSeeOther)
		return
	}

	err = database.ToggleCommentVote(userID, commentID, "like")
	if err != nil {
		log.Println("❌ erreur toggle comment like :", err)
		http.Redirect(w, r, "/echec", http.StatusSeeOther)
		return
	}

	http.Redirect(w, r, r.Referer(), http.StatusSeeOther)
}

func CommentDislikeHandler(w http.ResponseWriter, r *http.Request) {
	userID := GetCurrentUserID(r)
	commentIDStr := r.URL.Query().Get("id")
	commentID, err := strconv.Atoi(commentIDStr)
	if err != nil || userID == 0 {
		http.Redirect(w, r, "/echec", http.StatusSeeOther)
		return
	}

	err = database.ToggleCommentVote(userID, commentID, "dislike")
	if err != nil {
		log.Println("❌ erreur toggle comment dislike :", err)
		http.Redirect(w, r, "/echec", http.StatusSeeOther)
		return
	}

	http.Redirect(w, r, r.Referer(), http.StatusSeeOther)
}
func GetCurrentUserID(r *http.Request) int {
	cookie, err := r.Cookie("session")
	if err != nil {
		return 0
	}
	id, err := strconv.Atoi(cookie.Value)
	if err != nil {
		return 0
	}
	return id
}
func TimeHandlers(w http.ResponseWriter, r *http.Request) {
	// on recupere les posts triés par date
	posts, err := database.GetPostsSortedByDate()
	if err != nil {
		log.Println("❌ Erreur tri date :", err)
		http.Redirect(w, r, "/echec", http.StatusSeeOther)
		return
	}

	// structure pour date formatée
	var formatted []PostWithFormattedDate
	for _, post := range posts {
		formatted = append(formatted, PostWithFormattedDate{
			Post:          post,
			FormattedDate: post.Date.Format("02 Jan 2006 à 15:04"),
		})
	}

	data := struct {
		Posts []PostWithFormattedDate
	}{
		Posts: formatted,
	}

	tmpl, err := template.ParseFiles("templates/sorted.html") // crée une page si besoin
	if err != nil {
		log.Println("❌ Template trié manquant :", err)
		http.Redirect(w, r, "/echec", http.StatusSeeOther)
		return
	}

	tmpl.Execute(w, data)
}

func DeletePostHandler(w http.ResponseWriter, r *http.Request) {
	userID, _, err := GetCurrentUser(r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	postIDStr := r.URL.Query().Get("id")
	postID, err := strconv.Atoi(postIDStr)
	if err != nil {
		http.Redirect(w, r, "/echec", http.StatusSeeOther)
		return
	}

	// Vérifie que l'utilisateur est bien l'auteur du post
	post, err := database.GetPostByID(postID)
	if err != nil {
		http.Redirect(w, r, "/echec", http.StatusSeeOther)
		return
	}

	// Vérifie que l'utilisateur est bien l'auteur du post
	if post.Author != database.GetUsernameByID(userID) {
		http.Redirect(w, r, "/echec", http.StatusSeeOther)
		return
	}

	err = database.DeletePostByID(postID)
	if err != nil {
		http.Redirect(w, r, "/echec", http.StatusSeeOther)
		return
	}

	http.Redirect(w, r, "/home", http.StatusSeeOther)
}
