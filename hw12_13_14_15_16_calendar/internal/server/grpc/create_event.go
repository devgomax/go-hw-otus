package internalgrpc

import (
	"context"

	eventspb "github.com/devgomax/go-hw-otus/hw12_13_14_15_calendar/internal/pb/events"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// CreateEvent имплементация grpc метода CreateEvent.
func (i *Implementation) CreateEvent(ctx context.Context, req *eventspb.Event) (*eventspb.CreateEventResponse, error) {
	startsAt := req.StartsAt.AsTime()
	endsAt := req.EndsAt.AsTime()

	if err := i.app.CreateEvent(
		ctx,
		req.Id,
		req.Title,
		req.Description,
		req.UserId,
		&startsAt,
		&endsAt,
		req.NotifyInterval.AsDuration(),
	); err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to create event: %v", err)
	}

	return &eventspb.CreateEventResponse{}, nil
}
