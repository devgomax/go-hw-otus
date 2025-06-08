package adapters

import (
	"context"

	"github.com/devgomax/go-hw-otus/hw12_13_14_15_calendar/internal/app/sender"
	"github.com/devgomax/go-hw-otus/hw12_13_14_15_calendar/internal/pkg/clients/rabbitmq"
	"github.com/pkg/errors"
	amqp "github.com/rabbitmq/amqp091-go"
)

// Delivery адаптер RabbitMQ сообщения для сервиса планировщика.
type Delivery struct {
	amqp.Delivery
}

// GetBody возвращает тело RabbitMQ сообщения.
func (d Delivery) GetBody() []byte {
	return d.Body
}

// ClientRMQ адаптер RabbitMQ клиента для сервиса планировщика.
type ClientRMQ struct {
	*rabbitmq.Client
}

// NewRabbitMQClient конструктор адаптера RabbitMQ клиента для сервиса планировщика.
func NewRabbitMQClient(url, queue string) (*ClientRMQ, error) {
	client, err := rabbitmq.NewClient(url, queue)
	if err != nil {
		return nil, errors.Wrap(err, "[adaptermq::NewClient]")
	}

	return &ClientRMQ{client}, nil
}

// Consume адаптация Consume метода RabbitMQ клиента для сервиса планировщика.
func (c *ClientRMQ) Consume(ctx context.Context) (<-chan sender.IDeliveryMQ, error) {
	msgs, err := c.Client.Consume(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "[adaptermq::Consume]")
	}

	deliveries := make(chan sender.IDeliveryMQ)

	go func() {
		defer close(deliveries)
		for msg := range msgs {
			deliveries <- Delivery{msg}
		}
	}()

	return deliveries, nil
}
