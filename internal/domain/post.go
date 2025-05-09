package domain

import "time"

type Post struct {
	ID         string
	Title      string
	Content    string
	Author     string
	ImageURL   string
	Comments   []Comment
	CreatedAt  time.Time
	UserAvatar string
	AuthorID   string
}
type PostSummary struct {
	ID        string
	Title     string
	Author    string
	ImageURL  string
	CreatedAt time.Time
}
