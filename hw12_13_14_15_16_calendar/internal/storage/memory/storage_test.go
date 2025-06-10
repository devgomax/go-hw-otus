package memorystorage

import (
	"context"
	"slices"
	"strconv"
	"sync"
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

	t.Run("update event", func(t *testing.T) {
		t.Cleanup(func() {
			cleanup(repo)
		})

		err := repo.CreateEvent(context.Background(), events[0])
		require.NoError(t, err)

		eventUpd := &storage.Event{
			ID:             events[0].ID,
			Title:          "upd",
			StartsAt:       ptr(time.Now()),
			EndsAt:         ptr(time.Now()),
			Description:    "upd",
			UserID:         events[0].UserID,
			NotifyInterval: 10 * time.Second,
		}

		err = repo.UpdateEvent(context.Background(), eventUpd)
		require.NoError(t, err)

		require.Equal(t, eventUpd, repo.eventsByID[events[0].ID])
		require.Equal(t, eventUpd, repo.eventsByUser[events[0].UserID][0])
		require.Equal(t, eventUpd, repo.sortedEvents[0])
	})

	t.Run("read events", func(t *testing.T) {
		t.Cleanup(func() {
			cleanup(repo)
		})

		userID := "user"
		events := []*storage.Event{
			{
				Title:          "Today",
				StartsAt:       ptr(start),
				EndsAt:         ptr(start.Add(24 * time.Hour)),
				Description:    "Description",
				UserID:         userID,
				NotifyInterval: 15 * time.Second,
			},
			{
				Title:          "Title2",
				StartsAt:       ptr(start.Add(2 * 24 * time.Hour)),
				EndsAt:         ptr(start.Add(3 * 24 * time.Hour)),
				Description:    "Description",
				UserID:         userID,
				NotifyInterval: 15 * time.Second,
			},
			{
				Title:          "Title3",
				StartsAt:       ptr(start.Add(14 * 24 * time.Hour)),
				EndsAt:         ptr(start.Add(15 * 24 * time.Hour)),
				Description:    "Description",
				UserID:         userID,
				NotifyInterval: 15 * time.Second,
			},
		}

		for _, event := range events {
			err := repo.CreateEvent(context.Background(), event)
			require.NoError(t, err)
		}

		dailyEvents, err := repo.ReadDailyEvents(context.Background(), userID, start)
		require.NoError(t, err)
		require.Len(t, dailyEvents, 1)
		require.Equal(t, events[0], dailyEvents[0])

		dailyEvents, err = repo.ReadWeeklyEvents(context.Background(), userID, start)
		require.NoError(t, err)
		require.Equal(t, events[:2], dailyEvents)

		dailyEvents, err = repo.ReadMonthlyEvents(context.Background(), userID, start)
		require.NoError(t, err)
		require.Equal(t, events, dailyEvents)
	})
}

func TestStorageMultithreading(_ *testing.T) {
	repo := New()
	ctx := context.Background()

	userID := "user"

	start := time.Now()
	end := start.Add(24 * time.Hour)

	createFunc := func() {
		for range 1000 {
			_ = repo.CreateEvent(ctx, &storage.Event{
				StartsAt: &start,
				EndsAt:   &end,
				UserID:   userID,
			})
		}
	}

	readFunc := func() {
		for range 1000 {
			_, _ = repo.ReadDailyEvents(ctx, userID, start)
		}
	}

	wg := &sync.WaitGroup{}
	wg.Add(4)

	go func() {
		defer wg.Done()
		createFunc()
	}()

	go func() {
		defer wg.Done()
		createFunc()
	}()

	go func() {
		defer wg.Done()
		readFunc()
	}()

	go func() {
		defer wg.Done()
		readFunc()
	}()

	wg.Wait()
}
