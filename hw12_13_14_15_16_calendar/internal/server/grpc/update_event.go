package internalgrpc

import (
	"context"

	eventspb "github.com/devgomax/go-hw-otus/hw12_13_14_15_calendar/internal/pb/events"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// UpdateEvent имплементация grpc метода UpdateEvent.
func (i *Implementation) UpdateEvent(ctx context.Context, req *eventspb.Event) (*eventspb.UpdateEventResponse, error) {
	startsAt := req.StartsAt.AsTime()
	endsAt := req.EndsAt.AsTime()

	if err := i.app.UpdateEvent(
		ctx,
		req.Id,
		req.Title,
		req.Description,
		req.UserId,
		&startsAt,
		&endsAt,
		req.NotifyInterval.AsDuration(),
	); err != nil {
		return nil, status.Error(codes.Internal, "Failed to update event")
	}

	return &eventspb.UpdateEventResponse{}, nil
}
