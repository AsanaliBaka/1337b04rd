package domain

import (
	"1337b04rd/pkg"
	"time"
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
