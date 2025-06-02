package app

import (
	"context"

	"github.com/pkg/errors"
)

// DeleteEvent метод удаления события.
func (a *App) DeleteEvent(ctx context.Context, eventID string) error {
	err := a.repo.DeleteEvent(ctx, eventID)
	return errors.Wrapf(err, "[app::DeleteEvent]: failed to delete event by ID %q", eventID)
}
