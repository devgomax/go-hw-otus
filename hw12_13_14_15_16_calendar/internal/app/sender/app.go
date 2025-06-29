package sender

import (
	"context"
	"encoding/json"

	"github.com/devgomax/go-hw-otus/hw12_13_14_15_calendar/internal/app"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

// IConsumerMQ интерфейс консьюмера очереди сообщений.
type IConsumerMQ interface {
	Consume(ctx context.Context) (<-chan IDeliveryMQ, error)
	Close() error
}

// IDeliveryMQ интерфейс сообщения из очереди сообщений.
type IDeliveryMQ interface {
	Ack(multiple bool) error
	Reject(requeue bool) error
	Nack(multiple bool, requeue bool) error
	GetBody() []byte
}

// IRepository интерфейс БД.
type IRepository interface {
	SetEventsProcessedStatus(ctx context.Context, ids ...string) error
	Connect(ctx context.Context, dsn string) error
	Close()
}

// App структура приложения рассыльщика.
type App struct {
	repo     IRepository
	consumer IConsumerMQ
}

// NewApp конструктор приложения рассыльщика.
func NewApp(repo IRepository, consumer IConsumerMQ) *App {
	return &App{
		repo:     repo,
		consumer: consumer,
	}
}

// Run запускает непрерывное чтение сообщений из очереди.
func (a *App) Run(ctx context.Context) error {
	log.Info().Msg("[sender::Run]: start consuming...")

	msgs, err := a.consumer.Consume(ctx)
	if err != nil {
		return errors.Wrap(err, "[sender::Run]")
	}
	var notification app.EventNotification

	for msg := range msgs {
		body := msg.GetBody()
		if err = json.Unmarshal(body, &notification); err != nil {
			log.Error().Err(err).Bytes("notification", body).Msg(
				"[sender::Run]: failed to unmarshal amqp message")
			if err = msg.Reject(false); err != nil {
				log.Error().Err(err).Bytes("notification", body).Msg(
					"[sender::Run]: failed to reject amqp message")
			}

			continue
		}

		if err = a.repo.SetEventsProcessedStatus(ctx, notification.EventID); err != nil {
			log.Error().Err(err).Bytes("notification", body).Msg(
				"[sender::Run]: failed to update events DB statuses")
		}
	}

	return nil
}

// Stop закрывает amqp соединение.
func (a *App) Stop() error {
	if err := a.consumer.Close(); err != nil {
		return errors.Wrap(err, "[sender::Stop]")
	}

	log.Info().Msg("sender amqp connection gracefully shutdown")
	return nil
}
