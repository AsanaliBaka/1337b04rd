package domain

import (
	"time"

	"1337b04rd/pkg"
)

type Comment struct {
	ID        string
	Text      string
	Author    UserRef
	PostId    string
	CreatedAt time.Time
}

func NewComments(text string, author UserRef) *Comment {
	return &Comment{
		ID:        pkg.GeneratedId()(),
		Text:      text,
		Author:    author,
		CreatedAt: time.Now(),
	}
}
