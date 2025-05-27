package app

import (
	"context"

	"github.com/devgomax/go-hw-otus/hw12_13_14_15_calendar/internal/storage"
)

// App структура основного приложения.
type App struct {
	repo storage.Repository
}

// New конструктор основного приложения.
func New(repo storage.Repository) *App {
	return &App{
		repo: repo,
	}
}

// CreateEvent метод регистрации события.
func (a *App) CreateEvent(ctx context.Context, id, title string) error {
	// TODO
	return nil
	// return a.repo.CreateEvent(ctx, &storage.Event{ID: id, Title: title})
}

// TODO
