package models

import "time"

type Post struct {
    ID       int
    Author   string
    Content  string
    Category string
    Date     time.Time
}