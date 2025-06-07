package internalgrpc

import (
	"context"

	eventspb "github.com/devgomax/go-hw-otus/hw12_13_14_15_calendar/internal/pb/events"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// DeleteEvent имплементация grpc метода DeleteEvent.
func (i *Implementation) DeleteEvent(
	ctx context.Context,
	req *eventspb.DeleteEventRequest,
) (*eventspb.DeleteEventResponse, error) {
	if err := i.app.DeleteEvent(ctx, req.Id); err != nil {
		return nil, status.Error(codes.Internal, "Failed to delete event")
	}

	return &eventspb.DeleteEventResponse{}, nil
}
