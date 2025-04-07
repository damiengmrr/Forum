package database

import (
	"database/sql"
	"forum/models"
	"log"
)

// recupere les commentaires par post
func GetCommentsByPostID(postID int) ([]models.Comment, error) {
	query := "SELECT id, post_id, author, content, likes, dislikes, avatar, response_to FROM comments WHERE post_id = ?"
	rows, err := GetDatabase().Query(query, postID)
	if err != nil {
		log.Println("❌ erreur query comments:", err)
		return nil, err
	}
	defer rows.Close()

	var comments []models.Comment

	for rows.Next() {
		var c models.Comment
		var avatar sql.NullString

		err := rows.Scan(&c.ID, &c.PostID, &c.Author, &c.Content, &c.Likes, &c.Dislikes, &avatar, &c.ResponseTo)
		if err != nil {
			log.Println("❌ erreur scan commentaire:", err)
			continue
		}

		if avatar.Valid {
			c.Avatar = avatar.String
		} else {
			c.Avatar = "/static/image/default-avatar.png"
		}

		comments = append(comments, c)
	}

	return comments, nil
}

// insert un commentaire (reponse simple)
func InsertReply(postID int, content, author string) error {
	_, err := GetDatabase().Exec(`
		INSERT INTO comments (post_id, author, content, date, likes, dislikes)
			VALUES (?, ?, ?, datetime('now'), 0, 0)`, postID, author, content)
	if err != nil {
		log.Println("❌ erreur insert reply:", err)
	}
	return err
}

// increment like commentaire
func IncrementCommentLike(id int) error {
	_, err := GetDatabase().Exec("UPDATE comments SET likes = likes + 1 WHERE id = ?", id)
	if err != nil {
		log.Println("❌ erreur increment like comment:", err)
	}
	return err
}

// increment dislike commentaire
func IncrementCommentDislike(id int) error {
	_, err := GetDatabase().Exec("UPDATE comments SET dislikes = dislikes + 1 WHERE id = ?", id)
	if err != nil {
		log.Println("❌ erreur increment dislike comment:", err)
	}
	return err
}

// gestion votes commentaire (like/dislike exclusif)
func ToggleCommentVote(userID, commentID int, voteType string) error {
	db := GetDatabase()

	var current string
	err := db.QueryRow(`SELECT vote_type FROM votes_comments WHERE user_id = ? AND comment_id = ?`, userID, commentID).Scan(&current)

	switch {
	case err == sql.ErrNoRows:
		// pas encore vote -> insert
		_, err = db.Exec(`INSERT INTO votes_comments (user_id, comment_id, vote_type) VALUES (?, ?, ?)`, userID, commentID, voteType)
		if err != nil {
			log.Println("❌ erreur insert vote comment:", err)
			return err
		}

	case err != nil:
		log.Println("❌ erreur select vote comment:", err)
		return err

	case current == voteType:
		// meme vote -> delete
		_, err = db.Exec(`DELETE FROM votes_comments WHERE user_id = ? AND comment_id = ?`, userID, commentID)
		if err != nil {
			log.Println("❌ erreur delete vote comment:", err)
			return err
		}

	default:
		// vote different -> update
		_, err = db.Exec(`UPDATE votes_comments SET vote_type = ? WHERE user_id = ? AND comment_id = ?`, voteType, userID, commentID)
		if err != nil {
			log.Println("❌ erreur update vote comment:", err)
			return err
		}
	}

	// maj des compteurs apres le vote
	return updateCommentVoteCount(commentID)
}

// insert un commentaire avec possibilite de reponse
func InsertComment(postID int, author, content string, responseTo sql.NullInt64) error {
	_, err := GetDatabase().Exec(`
		INSERT INTO comments (post_id, author, content, date, likes, dislikes, response_to)
			VALUES (?, ?, ?, datetime('now'), 0, 0, ?)`, postID, author, content, responseTo)
	if err != nil {
		log.Println("❌ erreur insert comment:", err)
	}
	return err
}

// private function pour mettre a jour les compteurs like/dislike
func updateCommentVoteCount(commentID int) error {
	db := GetDatabase()

	var likes, dislikes int
	db.QueryRow(`SELECT COUNT(*) FROM votes_comments WHERE comment_id = ? AND vote_type = 'like'`, commentID).Scan(&likes)
	db.QueryRow(`SELECT COUNT(*) FROM votes_comments WHERE comment_id = ? AND vote_type = 'dislike'`, commentID).Scan(&dislikes)

	_, err := db.Exec(`UPDATE comments SET likes = ?, dislikes = ? WHERE id = ?`, likes, dislikes, commentID)
	if err != nil {
		log.Println("❌ erreur update compteur comment:", err)
	}
	return err
}
