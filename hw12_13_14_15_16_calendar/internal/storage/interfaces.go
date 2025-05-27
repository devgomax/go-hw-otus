package storage

import (
	"context"
	"time"
)

// Repository интерфейс хранилища, оперирующего событиями Event.
type Repository interface {
	Connect(ctx context.Context, dsn string) error
	Close()
	CreateEvent(ctx context.Context, event *Event) error
	UpdateEvent(ctx context.Context, event *Event) error
	DeleteEvent(ctx context.Context, eventID string) error
	ReadDailyEvents(ctx context.Context, userID string, date time.Time) ([]*Event, error)
	ReadWeeklyEvents(ctx context.Context, userID string, fromDate time.Time) ([]*Event, error)
	ReadMonthlyEvents(ctx context.Context, userID string, fromDate time.Time) ([]*Event, error)
}
