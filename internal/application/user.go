package application

import (
	"context"
	"time"

	"1337b04rd/internal/domain"
	"1337b04rd/pkg"
)

func (app *App) getUser(ctx context.Context, userID string) error {
	return nil
}

func (app *App) createUser(ctx context.Context) (*domain.User, error) {
	domain_user, err := app.avatarProvider.GetRandomAvatar()
	if err != nil {
		return nil, err
	}

	userID, err := pkg.GenerateUUID()
	if err != nil {
		return nil, err
	}
	newUser := &domain.User{
		ID:        userID,
		Username:  domain_user.Username,
		ImageURL:  domain_user.ImageURL,
		CreatedAt: time.Now(),
	}

	if err := app.repo.CreateUser(ctx, newUser); err != nil {
		return nil, err
	}

	return newUser, nil
}
