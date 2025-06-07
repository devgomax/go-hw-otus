package app

import (
	"github.com/devgomax/go-hw-otus/hw12_13_14_15_calendar/internal/storage"
)

// App структура основного приложения.
type App struct {
	repo storage.IRepository
}

// New конструктор основного приложения.
func New(repo storage.IRepository) *App {
	return &App{
		repo: repo,
	}
}
