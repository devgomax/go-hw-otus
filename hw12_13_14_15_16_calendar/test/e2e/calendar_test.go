package e2e

import (
	"context"
	"log"
	"os"
	"testing"
	"time"

	"github.com/devgomax/go-hw-otus/hw12_13_14_15_calendar/internal/app/scheduler"
	eventspb "github.com/devgomax/go-hw-otus/hw12_13_14_15_calendar/internal/pb/events"
	"github.com/devgomax/go-hw-otus/hw12_13_14_15_calendar/internal/pkg/clients/rabbitmq"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var (
	conn   *grpc.ClientConn
	client eventspb.EventsClient
	rabbit scheduler.IPublisherMQ
)

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	teardown()
	os.Exit(code)
}

func setup() {
	var err error

	conn, err = grpc.NewClient(":8091", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to create calendar client: %v", err)
	}

	client = eventspb.NewEventsClient(conn)

	rabbit, err = rabbitmq.NewClient("amqp://test:test@localhost:5673/", "events")
	if err != nil {
		log.Fatalf("failed to create amqp client: %v", err)
	}
}

func teardown() {
	if err := conn.Close(); err != nil {
		log.Printf("failed to close grpc connection: %v", err)
	}

	if err := rabbit.Close(); err != nil {
		log.Printf("failed to close amqp connection: %v", err)
	}
}

func TestCreateEvent(t *testing.T) {
	ctx := context.Background()
	now := time.Now().UTC()
	uid := uuid.New().String()
	event := &eventspb.Event{
		Id:             uid,
		Title:          uid,
		StartsAt:       timestamppb.New(now),
		EndsAt:         timestamppb.New(now.Add(1 * time.Hour)),
		Description:    "test1desc",
		UserId:         uid,
		NotifyInterval: durationpb.New(15 * time.Minute),
	}

	_, err := client.CreateEvent(ctx, event)
	require.NoError(t, err)

	resp, err := client.ReadDailyEvents(ctx, &eventspb.ReadDailyEventsRequest{
		UserId: uid,
		Date:   timestamppb.New(now),
	})
	require.NoError(t, err)

	require.Len(t, resp.Events, 1)
	require.Equal(t, event.Title, resp.Events[0].Title)

	event.Title = uuid.New().String()
	_, err = client.CreateEvent(ctx, event)
	require.Error(t, err)
}

func TestListingEvents(t *testing.T) {
	ctx := context.Background()
	now := time.Now().UTC()
	userID := uuid.New().String()

	dailyResp, err := client.ReadDailyEvents(ctx, &eventspb.ReadDailyEventsRequest{
		UserId: userID,
		Date:   timestamppb.New(now),
	})
	require.NoError(t, err)
	require.Empty(t, dailyResp.Events)

	weeklyResp, err := client.ReadWeeklyEvents(ctx, &eventspb.ReadWeeklyEventsRequest{
		UserId: userID,
		Date:   timestamppb.New(now),
	})
	require.NoError(t, err)
	require.Empty(t, weeklyResp.Events)

	monthlyResp, err := client.ReadMonthlyEvents(ctx, &eventspb.ReadMonthlyEventsRequest{
		UserId: userID,
		Date:   timestamppb.New(now),
	})
	require.NoError(t, err)
	require.Empty(t, monthlyResp.Events)

	getIDs := func(events []*eventspb.Event) []string {
		var ids []string
		for _, e := range events {
			ids = append(ids, e.Id)
		}
		return ids
	}

	events := make([]*eventspb.Event, 0, 3)

	for i := range 3 {
		mx := time.Duration(i)
		step := 6 * 24 * time.Hour
		uid := uuid.New().String()
		event := &eventspb.Event{
			Id:             uid,
			Title:          uid,
			StartsAt:       timestamppb.New(now.Add(step * mx)),
			EndsAt:         timestamppb.New(now.Add(step*mx + 1*time.Hour)),
			Description:    "test1desc",
			UserId:         userID,
			NotifyInterval: durationpb.New(15 * time.Minute),
		}
		events = append(events, event)
		_, err = client.CreateEvent(ctx, event)
		require.NoError(t, err)
	}

	dailyResp, err = client.ReadDailyEvents(ctx, &eventspb.ReadDailyEventsRequest{
		UserId: userID,
		Date:   timestamppb.New(now),
	})
	require.NoError(t, err)
	require.Len(t, dailyResp.Events, 1)
	require.Equal(t, events[0].Id, dailyResp.Events[0].Id)

	weeklyResp, err = client.ReadWeeklyEvents(ctx, &eventspb.ReadWeeklyEventsRequest{
		UserId: userID,
		Date:   timestamppb.New(now),
	})
	require.NoError(t, err)
	require.Len(t, weeklyResp.Events, 2)
	require.ElementsMatch(t, getIDs(events[:2]), getIDs(weeklyResp.Events))

	monthlyResp, err = client.ReadMonthlyEvents(ctx, &eventspb.ReadMonthlyEventsRequest{
		UserId: userID,
		Date:   timestamppb.New(now),
	})
	require.NoError(t, err)
	require.Len(t, monthlyResp.Events, len(events))
	require.ElementsMatch(t, getIDs(events), getIDs(monthlyResp.Events))
}

func TestReadDailyEventsCornerCases(t *testing.T) {
	ctx := context.Background()
	now := time.Now().UTC()
	todayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	todayEnd := todayStart.Add(24 * time.Hour)
	uid := uuid.New().String()

	tests := []struct {
		name  string
		start time.Time
		end   time.Time
		count int
	}{
		{
			name:  "event_end < search_start",
			start: todayStart.Add(-2 * time.Hour),
			end:   todayStart.Add(-1 * time.Hour),
		},
		{
			name:  "search_start = event_end",
			start: todayStart.Add(-1 * time.Hour),
			end:   todayStart,
		},
		{
			name:  "event_start < search_start < event_end",
			start: todayStart.Add(-1 * time.Hour),
			end:   todayStart.Add(1 * time.Hour),
			count: 1,
		},
		{
			name:  "search_start < event_start < event_end < search_end",
			start: todayStart.Add(1 * time.Hour),
			end:   todayEnd.Add(-1 * time.Hour),
			count: 1,
		},
		{
			name:  "event_start < search_end < event_end",
			start: todayEnd.Add(-1 * time.Hour),
			end:   todayEnd.Add(1 * time.Hour),
			count: 1,
		},
		{
			name:  "event_start = search_end",
			start: todayEnd,
			end:   todayEnd.Add(1 * time.Hour),
		},
		{
			name:  "search_end < event_start",
			start: todayEnd.Add(1 * time.Hour),
			end:   todayEnd.Add(2 * time.Hour),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			event := &eventspb.Event{
				Id:             uuid.New().String(),
				Title:          uid,
				StartsAt:       timestamppb.New(tc.start),
				EndsAt:         timestamppb.New(tc.end),
				Description:    "test1desc",
				UserId:         uid,
				NotifyInterval: durationpb.New(15 * time.Minute),
			}

			_, err := client.CreateEvent(ctx, event)
			require.NoError(t, err)

			resp, err := client.ReadDailyEvents(ctx, &eventspb.ReadDailyEventsRequest{
				UserId: uid,
				Date:   timestamppb.New(now),
			})
			require.NoError(t, err)
			require.Len(t, resp.Events, tc.count)
			if tc.count > 0 {
				require.Equal(t, event.Id, resp.Events[0].Id)
			}

			_, _ = client.DeleteEvent(ctx, &eventspb.DeleteEventRequest{Id: event.Id})
		})
	}
}

func TestReadWeeklyEventsCornerCases(t *testing.T) {
	ctx := context.Background()
	now := time.Now().UTC()
	todayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	weekEnd := todayStart.Add(7 * 24 * time.Hour)
	uid := uuid.New().String()

	tests := []struct {
		name  string
		start time.Time
		end   time.Time
		count int
	}{
		{
			name:  "event_end < search_start",
			start: todayStart.Add(-2 * time.Hour),
			end:   todayStart.Add(-1 * time.Hour),
		},
		{
			name:  "search_start = event_end",
			start: todayStart.Add(-1 * time.Hour),
			end:   todayStart,
		},
		{
			name:  "event_start < search_start < event_end",
			start: todayStart.Add(-1 * time.Hour),
			end:   todayStart.Add(1 * time.Hour),
			count: 1,
		},
		{
			name:  "search_start < event_start < event_end < search_end",
			start: todayStart.Add(1 * time.Hour),
			end:   weekEnd.Add(-1 * time.Hour),
			count: 1,
		},
		{
			name:  "event_start < search_end < event_end",
			start: weekEnd.Add(-1 * time.Hour),
			end:   weekEnd.Add(1 * time.Hour),
			count: 1,
		},
		{
			name:  "event_start = search_end",
			start: weekEnd,
			end:   weekEnd.Add(1 * time.Hour),
		},
		{
			name:  "search_end < event_start",
			start: weekEnd.Add(1 * time.Hour),
			end:   weekEnd.Add(2 * time.Hour),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			event := &eventspb.Event{
				Id:             uuid.New().String(),
				Title:          uid,
				StartsAt:       timestamppb.New(tc.start),
				EndsAt:         timestamppb.New(tc.end),
				Description:    "test1desc",
				UserId:         uid,
				NotifyInterval: durationpb.New(15 * time.Minute),
			}

			_, err := client.CreateEvent(ctx, event)
			require.NoError(t, err)

			resp, err := client.ReadWeeklyEvents(ctx, &eventspb.ReadWeeklyEventsRequest{
				UserId: uid,
				Date:   timestamppb.New(now),
			})
			require.NoError(t, err)
			require.Len(t, resp.Events, tc.count)
			if tc.count > 0 {
				require.Equal(t, event.Id, resp.Events[0].Id)
			}

			_, _ = client.DeleteEvent(ctx, &eventspb.DeleteEventRequest{Id: event.Id})
		})
	}
}

func TestReadMonthlyEventsCornerCases(t *testing.T) {
	ctx := context.Background()
	now := time.Now().UTC()
	todayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	monthEnd := time.Date(now.Year(), now.Month()+1, now.Day()+1, 0, 0, 0, 0, time.UTC)
	uid := uuid.New().String()

	tests := []struct {
		name  string
		start time.Time
		end   time.Time
		count int
	}{
		{
			name:  "event_end < search_start",
			start: todayStart.Add(-2 * time.Hour),
			end:   todayStart.Add(-1 * time.Hour),
		},
		{
			name:  "search_start = event_end",
			start: todayStart.Add(-1 * time.Hour),
			end:   todayStart,
		},
		{
			name:  "event_start < search_start < event_end",
			start: todayStart.Add(-1 * time.Hour),
			end:   todayStart.Add(1 * time.Hour),
			count: 1,
		},
		{
			name:  "search_start < event_start < event_end < search_end",
			start: todayStart.Add(1 * time.Hour),
			end:   monthEnd.Add(-1 * time.Hour),
			count: 1,
		},
		{
			name:  "event_start < search_end < event_end",
			start: monthEnd.Add(-1 * time.Hour),
			end:   monthEnd.Add(1 * time.Hour),
			count: 1,
		},
		{
			name:  "event_start = search_end",
			start: monthEnd,
			end:   monthEnd.Add(1 * time.Hour),
		},
		{
			name:  "search_end < event_start",
			start: monthEnd.Add(1 * time.Hour),
			end:   monthEnd.Add(2 * time.Hour),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			event := &eventspb.Event{
				Id:             uuid.New().String(),
				Title:          uid,
				StartsAt:       timestamppb.New(tc.start),
				EndsAt:         timestamppb.New(tc.end),
				Description:    "test1desc",
				UserId:         uid,
				NotifyInterval: durationpb.New(15 * time.Minute),
			}

			_, err := client.CreateEvent(ctx, event)
			require.NoError(t, err)

			resp, err := client.ReadMonthlyEvents(ctx, &eventspb.ReadMonthlyEventsRequest{
				UserId: uid,
				Date:   timestamppb.New(now),
			})
			require.NoError(t, err)
			require.Len(t, resp.Events, tc.count)
			if tc.count > 0 {
				require.Equal(t, event.Id, resp.Events[0].Id)
			}

			_, _ = client.DeleteEvent(ctx, &eventspb.DeleteEventRequest{Id: event.Id})
		})
	}
}
