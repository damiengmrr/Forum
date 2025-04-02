package database

import (
	"database/sql"
	"forum/models"
	"log"
)

func GetCommentsByPostID(postID int) ([]models.Comment, error) {
	db := GetDatabase()
	rows, err := db.Query("SELECT id, post_id, author, content, likes, dislikes, avatar, response_to FROM comments WHERE post_id = ?", postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []models.Comment

	for rows.Next() {
		var c models.Comment
		var avatar sql.NullString

		err := rows.Scan(&c.ID, &c.PostID, &c.Author, &c.Content, &c.Likes, &c.Dislikes, &avatar, &c.ResponseTo)
		if err != nil {
			log.Println("❌ erreur lecture commentaire :", err)
			continue
		}

		if avatar.Valid {
			c.Avatar = avatar.String
		} else {
			c.Avatar = "/static/image/default-avatar.png" // image par défaut
		}

		comments = append(comments, c)
	}

	return comments, nil
}

func InsertReply(postID int, content, author string) error {
	db := GetDatabase()
	stmt, err := db.Prepare("INSERT INTO comments (post_id, author, content, date, likes, dislikes) VALUES (?, ?, ?, datetime('now'), 0, 0)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(postID, author, content)
	return err
}
func IncrementCommentLike(id int) error { 
	db := GetDatabase()
	_, err := db.Exec("UPDATE comments SET likes = likes + 1 WHERE id = ?", id)
	return err
}

func IncrementCommentDislike(id int) error {
	db := GetDatabase()
	_, err := db.Exec("UPDATE comments SET dislikes = dislikes + 1 WHERE id = ?", id)
	return err
}

func ToggleCommentVote(userID, commentID int, voteType string) error {
	db := GetDatabase()

	var oldVote string
	err := db.QueryRow("SELECT vote_type FROM votes_comments WHERE user_id = ? AND comment_id = ?", userID, commentID).Scan(&oldVote)

	if err == sql.ErrNoRows {
		_, err = db.Exec("INSERT INTO votes_comments (user_id, comment_id, vote_type) VALUES (?, ?, ?)", userID, commentID, voteType)
	} else if err == nil {
		if oldVote == voteType {
			_, err = db.Exec("DELETE FROM votes_comments WHERE user_id = ? AND comment_id = ?", userID, commentID)
		} else {
			_, err = db.Exec("UPDATE votes_comments SET vote_type = ? WHERE user_id = ? AND comment_id = ?", voteType, userID, commentID)
		}
	} else {
		log.Println("❌ Erreur ToggleCommentVote :", err)
	}
	return err
}