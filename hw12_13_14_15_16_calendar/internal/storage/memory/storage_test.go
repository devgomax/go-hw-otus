package memorystorage

import (
	"context"
	"slices"
	"strconv"
	"testing"
	"time"

	"github.com/devgomax/go-hw-otus/hw12_13_14_15_calendar/internal/storage"
	"github.com/stretchr/testify/require"
)

func ptr[T any](val T) *T {
	return &val
}

func cleanup(repo *Repository) {
	repo.eventsByUser = make(map[string][]*storage.Event)
	repo.eventsByID = make(map[string]*storage.Event)
	repo.sortedEvents = make([]*storage.Event, 0)
}

func TestStorage(t *testing.T) {
	start := time.Now()

	events := make([]*storage.Event, 0, 5)

	for i := range 5 {
		str := strconv.Itoa(i)

		events = append(events, &storage.Event{
			Title:          "Title" + str,
			StartsAt:       ptr(start),
			EndsAt:         ptr(start.Add(10 * time.Second)),
			Description:    "Description" + str,
			UserID:         "user" + str,
			NotifyInterval: 15 * time.Second,
		})
	}

	repo := New()

	t.Run("new events successfully stored", func(t *testing.T) {
		t.Cleanup(func() {
			cleanup(repo)
		})

		for i, event := range events {
			err := repo.CreateEvent(context.Background(), event)
			require.NoError(t, err)
			require.Len(t, repo.sortedEvents, i+1)
			require.Len(t, repo.eventsByUser, i+1)
			require.Len(t, repo.eventsByID, i+1)
			require.Len(t, repo.eventsByUser[event.UserID], 1)
		}
	})

	t.Run("stored events are sorted", func(t *testing.T) {
		t.Cleanup(func() {
			cleanup(repo)
		})

		for _, event := range events {
			err := repo.CreateEvent(context.Background(), event)
			require.NoError(t, err)
		}

		eventsCopy := make([]*storage.Event, len(events))
		copy(eventsCopy, repo.sortedEvents)

		slices.SortFunc(eventsCopy, func(i, j *storage.Event) int {
			return i.StartsAt.Compare(*j.StartsAt)
		})

		require.Equal(t, repo.sortedEvents, eventsCopy)

	})

	t.Run("delete events", func(t *testing.T) {
		t.Cleanup(func() {
			cleanup(repo)
		})

		err := repo.CreateEvent(context.Background(), events[0])
		require.NoError(t, err)

		err = repo.DeleteEvent(context.Background(), events[0].ID)
		require.NoError(t, err)
		require.Empty(t, repo.eventsByID[events[0].ID])
		require.Empty(t, repo.sortedEvents)
		require.Empty(t, repo.eventsByUser[events[0].UserID])
	})
}
