package application

import (
	"context"
	"sync"
	"time"

	"1337b04rd/internal/ports/right"
)

type App struct {
	sync.Mutex
	userService    userService
	timers         map[string]*time.Timer
	ArchivePost    func(ctx context.Context, postID string)
	repo           right.DbPort
	avatarProvider right.AvatarProvider
	imageStorage   right.ImageStorage
}

func NewApp(pr right.DbPort, ar right.AvatarProvider, is right.ImageStorage, userService userService) *App {
	return &App{
		userService:    userService,
		timers:         make(map[string]*time.Timer),
		repo:           pr,
		avatarProvider: ar,
		imageStorage:   is,
	}
}

func (a *App) Timers() map[string]*time.Timer {
	return a.timers
}
