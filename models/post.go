package models

import (
	"database/sql"
	"time"
)

type Comment struct {
	ID         int
	PostID     int
	Author     string
	Avatar     string
	Content    string
	Likes      int
	Dislikes   int
	ResponseTo sql.NullInt64
	Response   *Comment
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
	Status     string
	Comments   []Comment
}
