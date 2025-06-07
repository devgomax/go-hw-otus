package memorystorage

import (
	"context"
	"slices"
	"sync"
	"time"

	"github.com/devgomax/go-hw-otus/hw12_13_14_15_calendar/internal/storage"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

// Repository модель БД типа in-memory.
type Repository struct {
	eventsByID   map[string]*storage.Event
	eventsByUser map[string][]*storage.Event
	sortedEvents []*storage.Event
	mu           sync.RWMutex
}

// New конструктор БД типа in-memory.
func New() *Repository {
	return &Repository{
		eventsByID:   make(map[string]*storage.Event),
		eventsByUser: make(map[string][]*storage.Event),
	}
}

// Connect открывает соединение с БД.
func (r *Repository) Connect(_ context.Context, _ string) error { return nil }

// Close закрывает соединение с БД.
func (r *Repository) Close() {}

// CreateEvent сохраняет событие в БД.
func (r *Repository) CreateEvent(_ context.Context, event *storage.Event) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.eventsByID[event.ID]; exists {
		return errors.Errorf("[memorystorage::CreateEvent]: event with ID %s already exists", event.ID)
	}

	event.ID = uuid.New().String() // имитируем поведение "UUID PRIMARY KEY" как в postgres
	r.eventsByID[event.ID] = event
	r.eventsByUser[event.UserID] = append(r.eventsByUser[event.UserID], event)

	r.sortedEvents = append(r.sortedEvents, event)
	slices.SortFunc(r.sortedEvents, func(i, j *storage.Event) int {
		return i.StartsAt.Compare(*j.StartsAt)
	})

	return nil
}

// UpdateEvent обновляет событие в БД.
func (r *Repository) UpdateEvent(_ context.Context, event *storage.Event) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.eventsByID[event.ID]; !exists {
		return errors.Errorf("[memorystorage::UpdateEvent]: event with ID %s does not exist", event.ID)
	}

	r.eventsByID[event.ID] = event

	for i, e := range r.eventsByUser[event.UserID] {
		if e.ID == event.ID {
			r.eventsByUser[event.UserID][i] = event
			break
		}
	}

	for i, e := range r.sortedEvents {
		if e.ID == event.ID {
			r.sortedEvents[i] = event
			break
		}
	}

	slices.SortFunc(r.sortedEvents, func(i, j *storage.Event) int {
		return i.StartsAt.Compare(*j.StartsAt)
	})

	return nil
}

// DeleteEvent удаляет событие из БД.
func (r *Repository) DeleteEvent(_ context.Context, eventID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	event, exists := r.eventsByID[eventID]
	if !exists {
		return errors.Errorf("[memorystorage::UpdateEvent]: event with ID %s does not exist", eventID)
	}

	delete(r.eventsByID, eventID)

	userEvents := r.eventsByUser[event.UserID]
	for i, e := range userEvents {
		if e.ID == eventID {
			r.eventsByUser[event.UserID] = slices.Delete(r.eventsByUser[event.UserID], i, i+1)
			break
		}
	}

	for i, e := range r.sortedEvents {
		if e.ID == eventID {
			r.sortedEvents = slices.Delete(r.sortedEvents, i, i+1)
			break
		}
	}

	return nil
}

// ReadDailyEvents читает события за указанную дату.
func (r *Repository) ReadDailyEvents(_ context.Context, userID string, date time.Time) ([]*storage.Event, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	start := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	end := start.Add(24 * time.Hour)

	var result []*storage.Event
	for _, event := range r.sortedEvents {
		if event.UserID == userID && event.StartsAt.Before(end) && event.EndsAt.After(start) {
			result = append(result, event)
		}
	}

	return result, nil
}

// ReadWeeklyEvents читает события за неделю, начиная с указанной даты.
func (r *Repository) ReadWeeklyEvents(_ context.Context, userID string, fromDate time.Time) ([]*storage.Event, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	start := time.Date(fromDate.Year(), fromDate.Month(), fromDate.Day(), 0, 0, 0, 0, fromDate.Location())
	end := start.Add(7 * 24 * time.Hour)

	var result []*storage.Event
	for _, event := range r.sortedEvents {
		if event.UserID == userID && event.StartsAt.Before(end) && event.EndsAt.After(start) {
			result = append(result, event)
		}
	}

	return result, nil
}

// ReadMonthlyEvents читает события за месяц, начиная с указанной даты.
func (r *Repository) ReadMonthlyEvents(_ context.Context, userID string, fromDate time.Time) ([]*storage.Event, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	start := time.Date(fromDate.Year(), fromDate.Month(), fromDate.Day(), 0, 0, 0, 0, fromDate.Location())
	end := time.Date(fromDate.Year(), fromDate.Month()+1, fromDate.Day()+1, 0, 0, 0, 0, fromDate.Location())

	var result []*storage.Event
	for _, event := range r.sortedEvents {
		if event.UserID == userID && event.StartsAt.Before(end) && event.EndsAt.After(start) {
			result = append(result, event)
		}
	}

	return result, nil
}
