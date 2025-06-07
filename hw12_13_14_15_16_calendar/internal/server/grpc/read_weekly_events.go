package internalgrpc

import (
	"context"

	eventspb "github.com/devgomax/go-hw-otus/hw12_13_14_15_calendar/internal/pb/events"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// ReadWeeklyEvents имплементация grpc метода ReadWeeklyEvents.
func (i *Implementation) ReadWeeklyEvents(
	ctx context.Context,
	req *eventspb.ReadWeeklyEventsRequest,
) (*eventspb.ReadWeeklyEventsResponse, error) {
	events, err := i.app.ReadWeeklyEvents(ctx, req.UserId, req.Date.AsTime())
	if err != nil {
		return nil, status.Error(codes.Internal, "Failed to get weekly events")
	}

	var response eventspb.ReadWeeklyEventsResponse

	for _, event := range events {
		response.Events = append(response.Events, &eventspb.Event{
			Id:             event.ID,
			Title:          event.Title,
			StartsAt:       timestamppb.New(*event.StartsAt),
			EndsAt:         timestamppb.New(*event.EndsAt),
			Description:    event.Description,
			UserId:         event.UserID,
			NotifyInterval: durationpb.New(event.NotifyInterval),
		})
	}

	return &response, nil
}
