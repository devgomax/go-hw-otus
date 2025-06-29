package scheduler

import (
	"context"
	"encoding/json"
	"time"

	"github.com/devgomax/go-hw-otus/hw12_13_14_15_calendar/internal/app"
	"github.com/devgomax/go-hw-otus/hw12_13_14_15_calendar/internal/storage"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

// IPublisherMQ интерфейс отправителя очереди сообщений.
type IPublisherMQ interface {
	Publish(ctx context.Context, msg []byte) error
	Close() error
}

// IRepository интерфейс БД.
type IRepository interface {
	ReadEventsToNotify(ctx context.Context) ([]*storage.Event, error)
	Connect(ctx context.Context, dsn string) error
	Close()
}

// App структура приложения планировщика.
type App struct {
	repo      IRepository
	publisher IPublisherMQ
}

// NewApp конструктор приложения планировщика.
func NewApp(repo IRepository, publisher IPublisherMQ) *App {
	return &App{
		repo:      repo,
		publisher: publisher,
	}
}

// Run запускает периодическое сканирование БД и отправку уведомлений о событиях в очередь.
func (a *App) Run(ctx context.Context, tickd time.Duration) error {
	ticker := time.NewTicker(tickd)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
			events, err := a.repo.ReadEventsToNotify(ctx)
			if err != nil {
				return errors.Wrap(err, "[scheduler::Run]: failed to read events from DB")
			}

			for _, event := range events {
				notification := app.EventNotification{
					EventID:    event.ID,
					EventTitle: event.Title,
					EventDate:  *event.StartsAt,
					UserID:     event.UserID,
				}

				data, inErr := json.Marshal(notification)
				if inErr != nil {
					return errors.Wrap(inErr, "[scheduler::Run]: can't marshal notification")
				}

				if err = a.publisher.Publish(ctx, data); err != nil {
					return errors.Wrap(err, "[scheduler::Run]: failed to publish amqp message")
				}
			}
		}
	}
}

// Stop закрывает amqp соединение.
func (a *App) Stop() error {
	if err := a.publisher.Close(); err != nil {
		return errors.Wrap(err, "[scheduler::Stop]")
	}

	log.Info().Msg("scheduler amqp connection gracefully shutdown")
	return nil
}
