package domain

import (
	"time"

	"1337b04rd/pkg"
)

type UserRef struct {
	SessionID string
	AvatarURL string
	Name      string
	ExpiresAt time.Time
}

func NewUser(avatar, name string) *UserRef {
	return &UserRef{
		SessionID: pkg.GeneratedId()(),
		AvatarURL: avatar,
		Name:      name,
	}
}
