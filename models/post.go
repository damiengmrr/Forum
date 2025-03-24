package models

import "time"

// Structure Post mise à jour
type Post struct {
	ID         int
	Author     string
	Title      string
	Content    string
	Categories []string
	Date       time.Time
	ImagePath  string
	Likes      int
	Dislikes   int
	Status     string // "published", "draft", "abandoned"
}
