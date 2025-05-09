package right

import "1337b04rd/internal/domain"

type AvatarProvider interface {
	GetRandomAvatar() (*domain.User, error)
	GetRandomAvatarByID(id int) (*domain.User, error)
}
