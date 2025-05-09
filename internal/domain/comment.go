package domain

import "time"

type Comment struct {
	ID         string
	Author     string
	Content    string
	AvaterLink string
	ParentID   string
	Replies    []Comment
	CreatedAt  time.Time
}
