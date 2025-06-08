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

// App структура приложения рассыльщика.
type App struct {
	consumer IConsumerMQ
}

// NewApp конструктор приложения рассыльщика.
func NewApp(consumer IConsumerMQ) *App {
	return &App{
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

		log.Info().RawJSON("notification", body).Msg("notification sent successfully")
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
