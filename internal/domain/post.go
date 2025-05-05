package domain

import "time"

type Post struct {
	ID        string
	Title     string
	Content   string
	ImageURL  string
	User      User
	CreatedAt time.Time
	Comments  []Comment
}

type Comment struct {
	ID        string
	PostID    string
	ParentID  *string // nil если ответ на пост
	Content   string
	User      User
	CreatedAt time.Time
}

type User struct {
	SessionID string
	Name      string
	AvatarURL string
}
