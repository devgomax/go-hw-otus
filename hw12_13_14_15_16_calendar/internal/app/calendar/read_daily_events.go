package calendar

import (
	"context"
	"time"

	"github.com/devgomax/go-hw-otus/hw12_13_14_15_calendar/internal/storage"
	"github.com/pkg/errors"
)

// ReadDailyEvents метод получения событий за определенную дату.
func (a *App) ReadDailyEvents(ctx context.Context, userID string, date time.Time) ([]*storage.Event, error) {
	events, err := a.repo.ReadDailyEvents(ctx, userID, date)

	return events, errors.Wrapf(err,
		"[app::ReadDailyEvents]: failed to get daily events by user %q and date %v", userID, date)
}
