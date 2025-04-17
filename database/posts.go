package database

import (
	"database/sql"
	"forum/models"
	"log"
	"strings"
	"time"
)

// insert un post dans la bdd
func InsertPost(post models.Post) error {
	_, err := GetDatabase().Exec(`
		INSERT INTO posts(author, title, content, date, image_path, categories, likes, dislikes)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		post.Author, post.Title, post.Content,
		post.Date.Format("2006-01-02 15:04:05"), post.ImagePath,
		strings.Join(post.Categories, ","), post.Likes, post.Dislikes,
	)
	if err != nil {
		log.Println("❌ erreur insert post:", err)
	}
	return err
}

// recup tous les posts
func GetAllPosts() ([]models.Post, error) {
	return fetchPosts("SELECT id, author, title, content, date, image_path, categories, likes, dislikes FROM posts ORDER BY date DESC")
}

// recup les posts par categorie
func GetPostsByCategory(category string) ([]models.Post, error) {
	query := "SELECT id, author, title, content, date, image_path, categories, likes, dislikes FROM posts WHERE categories LIKE ? ORDER BY date DESC"
	return fetchPosts(query, "%"+category+"%")
}

// recup un post par id
func GetPostByID(id int) (models.Post, error) {
	query := "SELECT id, author, title, content, date, image_path, categories, likes, dislikes FROM posts WHERE id = ?"
	posts, err := fetchPosts(query, id)
	if err != nil || len(posts) == 0 {
		if err != nil {
			log.Println("❌ erreur fetch post by id:", err)
		}
		return models.Post{}, sql.ErrNoRows
	}
	return posts[0], nil
}

// private function pour factoriser les fetch de posts
func fetchPosts(query string, args ...interface{}) ([]models.Post, error) {
	rows, err := GetDatabase().Query(query, args...)
	if err != nil {
		log.Println("❌ erreur query posts:", err)
		return nil, err
	}
	defer rows.Close()

	var posts []models.Post

	for rows.Next() {
		var post models.Post
		var dateStr, catString string

		err := rows.Scan(&post.ID, &post.Author, &post.Title, &post.Content,
			&dateStr, &post.ImagePath, &catString, &post.Likes, &post.Dislikes)
		if err != nil {
			log.Println("❌ erreur scan post:", err)
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
		//fmt.Print(err)
	}
	defer rows.Close()

	var posts []models.Post
	for rows.Next() {
		var post models.Post
		var dateStr, catString string

		err := rows.Scan(&post.ID, &post.Author, &post.Title, &post.Content, &dateStr, &post.ImagePath, &catString, &post.Likes, &post.Dislikes)
		if err != nil {
			log.Println("❌ erreur parse date:", err)
		}

		post.Categories = strings.Split(catString, ",")
		posts = append(posts, post)
	}

	return posts, nil
}

// tri par date
func GetPostsSortedByDate() ([]models.Post, error) {
	return GetAllPosts() // car deja trie par date dans GetAllPosts
}

// gestion votes posts (like/dislike exclusif)
func TogglePostVote(userID, postID int, voteType string) error {
	db := GetDatabase()

	var oldVote string
	err := db.QueryRow(`SELECT vote_type FROM votes_posts WHERE user_id = ? AND post_id = ?`, userID, postID).Scan(&oldVote)

	switch {
	case err == sql.ErrNoRows:
		// pas encore vote -> insert
		_, err := db.Exec(`INSERT INTO votes_posts (user_id, post_id, vote_type) VALUES (?, ?, ?)`, userID, postID, voteType)
		if err != nil {
			log.Println("❌ erreur insert vote post:", err)
			return err
		}
		return updatePostVoteCount(postID, voteType, 1)

	case err != nil:
		log.Println("❌ erreur select vote post:", err)
		return err

	case oldVote == voteType:
		// meme vote -> delete
		_, err := db.Exec(`DELETE FROM votes_posts WHERE user_id = ? AND post_id = ?`, userID, postID)
		if err != nil {
			log.Println("❌ erreur delete vote post:", err)
			return err
		}
		return updatePostVoteCount(postID, voteType, -1)

	default:
		// vote different -> update
		_, err := db.Exec(`UPDATE votes_posts SET vote_type = ? WHERE user_id = ? AND post_id = ?`, voteType, userID, postID)
		if err != nil {
			log.Println("❌ erreur update vote post:", err)
			return err
		}
		// on inverse les compteurs
		if voteType == "like" {
			_ = updatePostVoteCount(postID, "like", 1)
			_ = updatePostVoteCount(postID, "dislike", -1)
		} else {
			_ = updatePostVoteCount(postID, "dislike", 1)
			_ = updatePostVoteCount(postID, "like", -1)
		}
		return nil
	}
}

// private function pour update les compteurs likes/dislikes
func updatePostVoteCount(postID int, voteType string, delta int) error {
	db := GetDatabase()
	query := fmt.Sprintf("UPDATE posts SET %s = %s + ? WHERE id = ?", voteType, voteType)
	_, err := db.Exec(query, delta, postID)
	if err != nil {
		log.Println("❌ erreur update compteur vote post:", err)
	}
	return err
}

func DeletePostByID(postID int) error {
	db := GetDatabase()
	_, err := db.Exec("DELETE FROM posts WHERE id = ?", postID)
	return err
}
