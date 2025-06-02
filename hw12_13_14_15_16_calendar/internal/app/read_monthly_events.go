package app

import (
	"context"
	"time"

	"github.com/devgomax/go-hw-otus/hw12_13_14_15_calendar/internal/storage"
	"github.com/pkg/errors"
)

// ReadMonthlyEvents метод получения событий за месяц, начиная с определенной даты.
func (a *App) ReadMonthlyEvents(ctx context.Context, userID string, date time.Time) ([]*storage.Event, error) {
	events, err := a.repo.ReadMonthlyEvents(ctx, userID, date)

	return events, errors.Wrapf(err,
		"[app::ReadMonthlyEvents]: failed to get weekly events by user %q and date %v", userID, date)
}
