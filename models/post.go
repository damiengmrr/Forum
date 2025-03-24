package models

import "time"

type Comment struct {
	ID       int
	Author   string
	Avatar   string
	Content  string
	Likes    int
	Dislikes int
	Response *Comment // une seule réponse max
}

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
	Comments   []Comment
}