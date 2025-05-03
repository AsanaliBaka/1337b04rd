package domain

import (
	"time"

	"1337b04rd/pkg"
)

type Comment struct {
	ID        string
	Text      string
	Author    UserRef
	ImageURL  string
	CreatedAt time.Time
}

func NewComments(text, image string, author UserRef) *Comment {
	return &Comment{
		ID:        pkg.GeneratedId()(),
		Text:      text,
		Author:    author,
		ImageURL:  image,
		CreatedAt: time.Now(),
	}
}
