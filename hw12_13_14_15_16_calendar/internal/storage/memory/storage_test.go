package memorystorage

import (
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
			ID:             str,
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
			err := repo.CreateEvent(t.Context(), event)
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
			err := repo.CreateEvent(t.Context(), event)
			require.NoError(t, err)
		}

		for i, event := range repo.sortedEvents {
			require.Equal(t, strconv.Itoa(i), event.ID)
		}
	})

	t.Run("delete events", func(t *testing.T) {
		t.Cleanup(func() {
			cleanup(repo)
		})

		err := repo.CreateEvent(t.Context(), events[0])
		require.NoError(t, err)

		err = repo.DeleteEvent(t.Context(), events[0].ID)
		require.NoError(t, err)
		require.Empty(t, repo.eventsByID[events[0].ID])
		require.Empty(t, repo.sortedEvents)
		require.Empty(t, repo.eventsByUser[events[0].UserID])
	})
}
