package application

import (
	"errors"
	"math/rand"
)

type userService struct{}

func NewUser() *userService {
	return &userService{}
}

func (user *userService) GenUserID(userCount, maxID int) (int, error) {
	if maxID == 0 {
		return 0, errors.New("maxID is zero")
	}

	id := userCount + 1
	if userCount > maxID {
		id = rand.Intn(maxID)
	}

	return id, nil
}
