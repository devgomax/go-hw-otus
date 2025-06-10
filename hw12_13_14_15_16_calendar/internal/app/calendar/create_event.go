package calendar

import (
	"context"
	"time"

	"github.com/devgomax/go-hw-otus/hw12_13_14_15_calendar/internal/storage"
	"github.com/pkg/errors"
)

// CreateEvent метод регистрации события.
func (a *App) CreateEvent(
	ctx context.Context,
	title string,
	description string,
	userID string,
	startsAt *time.Time,
	endAt *time.Time,
	notifyInterval time.Duration,
) error {
	event := storage.Event{
		Title:          title,
		StartsAt:       startsAt,
		EndsAt:         endAt,
		Description:    description,
		UserID:         userID,
		NotifyInterval: notifyInterval,
	}

	err := a.repo.CreateEvent(ctx, &event)

	return errors.Wrap(err, "[app::CreateEvent]: failed to create event")
}
