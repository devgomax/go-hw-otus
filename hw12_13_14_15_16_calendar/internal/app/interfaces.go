package app

import (
	"context"
	"time"

	"github.com/devgomax/go-hw-otus/hw12_13_14_15_calendar/internal/storage"
)

// IApp основной API интерфейс.
type IApp interface {
	CreateEvent(ctx context.Context, title, description, userID string, startsAt, endAt *time.Time, notifyInterval time.Duration) error //nolint:lll
	DeleteEvent(ctx context.Context, eventID string) error
	ReadDailyEvents(ctx context.Context, userID string, date time.Time) ([]*storage.Event, error)
	ReadWeeklyEvents(ctx context.Context, userID string, date time.Time) ([]*storage.Event, error)
	ReadMonthlyEvents(ctx context.Context, userID string, date time.Time) ([]*storage.Event, error)
	UpdateEvent(ctx context.Context, id, title, description, userID string, startsAt, endAt *time.Time, notifyInterval time.Duration) error //nolint:lll
}
