package app

import (
	"context"
	"time"

	"github.com/devgomax/go-hw-otus/hw12_13_14_15_calendar/internal/storage"
	"github.com/pkg/errors"
)

// UpdateEvent метод обновления события.
func (a *App) UpdateEvent(
	ctx context.Context,
	id string,
	title string,
	description string,
	userID string,
	startsAt *time.Time,
	endAt *time.Time,
	notifyInterval time.Duration,
) error {
	event := storage.Event{
		ID:             id,
		Title:          title,
		StartsAt:       startsAt,
		EndsAt:         endAt,
		Description:    description,
		UserID:         userID,
		NotifyInterval: notifyInterval,
	}

	err := a.repo.UpdateEvent(ctx, &event)

	return errors.Wrapf(err, "[app::UpdateEvent]: failed to update event with ID %q", id)
}
