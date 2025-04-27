package domain

import (
	"time"
)

type Post struct {
	ID         string //id поста
	Title      string
	Content    string
	ImageURL   string
	Author     UserRef
	CreatedAt  time.Time
	Expires    time.Time
	IsArchived bool
}

func NewPost(title, content string, author UserRef) *Post {
	return &Post{
		ID:        generatedId(),
		Title:     title,
		Content:   content,
		Author:    author,
		CreatedAt: time.Now(),
		Expires:   time.Now().Add(10 * time.Minute),
	}
}
