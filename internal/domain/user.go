package domain

import "time"

type User struct {
	ID        string
	Username  string
	ImageURL  string
	CreatedAt time.Time
}
