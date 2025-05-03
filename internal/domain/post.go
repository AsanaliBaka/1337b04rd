package domain

import (
	"time"

	"1337b04rd/pkg"
)

type Post struct {
	ID         string
	Title      string
	Content    string
	ImageURL   string
	Author     UserRef
	CreatedAt  time.Time
	UpdatedAt  time.Time
	IsArchived bool
}

func NewPost(title, content string, author UserRef) *Post {
	return &Post{
		ID:        pkg.GeneratedId()(),
		Title:     title,
		Content:   content,
		Author:    author,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now().Add(10 * time.Minute),
	}
}
