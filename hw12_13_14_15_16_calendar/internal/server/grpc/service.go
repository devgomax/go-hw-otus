package internalgrpc

import (
	"github.com/devgomax/go-hw-otus/hw12_13_14_15_calendar/internal/app"
	eventspb "github.com/devgomax/go-hw-otus/hw12_13_14_15_calendar/internal/pb/events"
)

// Implementation имплементация grpc-сервера Events.
type Implementation struct {
	eventspb.UnimplementedEventsServer

	app *app.App
}

// NewEventsServer конструктор grpc-сервера Events.
func NewEventsServer(app *app.App) *Implementation {
	return &Implementation{
		app: app,
	}
}
