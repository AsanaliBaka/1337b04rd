package domain

import (
	. "1337b04rd/pkg"
)

type UserRef struct {
	SessionID string
	AvatarURL string
	Name      string
}

func NewUser(avatar, name string) *UserRef {
	return &UserRef{
		SessionID: GeneratedId()(),
		AvatarURL: avatar,
		Name:      name,
	}
}
