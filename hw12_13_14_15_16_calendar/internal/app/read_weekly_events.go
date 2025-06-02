package app

import (
	"context"
	"time"

	"github.com/devgomax/go-hw-otus/hw12_13_14_15_calendar/internal/storage"
	"github.com/pkg/errors"
)

// ReadWeeklyEvents метод получения событий за неделю, начиная с определенной даты.
func (a *App) ReadWeeklyEvents(ctx context.Context, userID string, date time.Time) ([]*storage.Event, error) {
	events, err := a.repo.ReadWeeklyEvents(ctx, userID, date)

	return events, errors.Wrapf(err,
		"[app::ReadWeeklyEvents]: failed to get weekly events by user %q and date %v", userID, date)
}
