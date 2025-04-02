package database

import (
	"database/sql"
	"forum/models"
	"log"
	"strings"
	"time"
)

// ajoute un post dans la bdd
func InsertPost(post models.Post) error {
	stmt, err := GetDatabase().Prepare("INSERT INTO posts(author, title, content, date, image_path, categories, likes, dislikes) VALUES (?, ?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(post.Author, post.Title, post.Content, post.Date.Format("2006-01-02 15:04:05"), post.ImagePath, strings.Join(post.Categories, ","), post.Likes, post.Dislikes)
	return err
}

// recupere tous les posts
func GetAllPosts() ([]models.Post, error) {
	rows, err := GetDatabase().Query("SELECT id, author, title, content, date, image_path, categories, likes, dislikes FROM posts ORDER BY date DESC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []models.Post
	for rows.Next() {
		var post models.Post
		var dateStr, catString string

		err := rows.Scan(&post.ID, &post.Author, &post.Title, &post.Content, &dateStr, &post.ImagePath, &catString, &post.Likes, &post.Dislikes)
		if err != nil {
			log.Println("erreur scan :", err)
			continue
		}

		post.Date, _ = time.Parse("2006-01-02 15:04:05", dateStr)
		post.Categories = strings.Split(catString, ",")

		posts = append(posts, post)
	}

	return posts, nil
}

// recupere les posts par categorie
func GetPostsByCategory(category string) ([]models.Post, error) {
	rows, err := GetDatabase().Query("SELECT id, author, title, content, date, image_path, categories, likes, dislikes FROM posts WHERE categories LIKE ? ORDER BY date DESC", "%"+category+"%")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []models.Post
	for rows.Next() {
		var post models.Post
		var dateStr, catString string

		err := rows.Scan(&post.ID, &post.Author, &post.Title, &post.Content, &dateStr, &post.ImagePath, &catString, &post.Likes, &post.Dislikes)
		if err != nil {
			log.Println("erreur scan :", err)
			continue
		}

		post.Date, _ = time.Parse("2006-01-02 15:04:05", dateStr)
		post.Categories = strings.Split(catString, ",")

		posts = append(posts, post)
	}

	return posts, nil
}

// recupere un post par id
func GetPostByID(id int) (models.Post, error) {
	var post models.Post
	var dateStr, catStr string

	row := GetDatabase().QueryRow(`
		SELECT id, author, title, content, date, image_path, categories, likes, dislikes
		FROM posts WHERE id = ?`, id)

	err := row.Scan(&post.ID, &post.Author, &post.Title, &post.Content,
		&dateStr, &post.ImagePath, &catStr, &post.Likes, &post.Dislikes)
	if err != nil {
		return post, err
	}

	// formatage de la date
	post.Date, _ = time.Parse("2006-01-02 15:04:05", dateStr)
	post.Categories = strings.Split(catStr, ",")

	return post, nil
}

func IncrementLike(id int) error {
	db := GetDatabase()
	_, err := db.Exec("UPDATE posts SET likes = likes + 1 WHERE id = ?", id)
	return err
}

func IncrementDislike(id int) error {
	db := GetDatabase()
	_, err := db.Exec("UPDATE posts SET dislikes = dislikes + 1 WHERE id = ?", id)
	return err
}
func GetPostsSortedByDate() ([]models.Post, error) {
	db := GetDatabase()
	rows, err := db.Query("SELECT id, author, title, content, date, image_path, categories, likes, dislikes FROM posts ORDER BY date DESC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []models.Post
	for rows.Next() {
		var post models.Post
		var dateStr, catString string

		err := rows.Scan(&post.ID, &post.Author, &post.Title, &post.Content, &dateStr, &post.ImagePath, &catString, &post.Likes, &post.Dislikes)
		if err != nil {
			log.Println("scan post erreur :", err)
			continue
		}

		post.Date, _ = time.Parse("2006-01-02 15:04:05", dateStr)
		post.Categories = strings.Split(catString, ",")
		posts = append(posts, post)
	}

	return posts, nil
}

func TogglePostVote(userID, postID int, voteType string) error {
	db := GetDatabase()

	var oldVote string
	err := db.QueryRow("SELECT vote_type FROM votes_posts WHERE user_id = ? AND post_id = ?", userID, postID).Scan(&oldVote)

	if err == sql.ErrNoRows {
		_, err = db.Exec("INSERT INTO votes_posts (user_id, post_id, vote_type) VALUES (?, ?, ?)", userID, postID, voteType)
	} else if err == nil {
		if oldVote == voteType {
			_, err = db.Exec("DELETE FROM votes_posts WHERE user_id = ? AND post_id = ?", userID, postID)
		} else {
			_, err = db.Exec("UPDATE votes_posts SET vote_type = ? WHERE user_id = ? AND post_id = ?", voteType, userID, postID)
		}
	} else {
		log.Println("❌ Erreur TogglePostVote :", err)
	}
	return err
}
