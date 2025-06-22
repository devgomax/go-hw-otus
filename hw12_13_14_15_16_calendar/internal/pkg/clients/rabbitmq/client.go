package rabbitmq

import (
	"context"

	"github.com/pkg/errors"
	amqp "github.com/rabbitmq/amqp091-go"
)

// Client основной RabbitMQ клиент.
type Client struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	queue   string
}

// NewClient конструктор RabbitMQ клиента.
func NewClient(url, queue string) (*Client, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, errors.Wrap(err, "[rabbitmq::NewClient]: failed to establish amqp connection")
	}

	channel, err := conn.Channel()
	if err != nil {
		return nil, errors.Wrap(err, "[rabbitmq::NewClient]: failed to open amqp channel")
	}

	if _, err = channel.QueueDeclare(queue, true, false, false, false, nil); err != nil {
		return nil, errors.Wrapf(err, "[rabbitmq::NewClient]: failed to declare amqp queue %q", queue)
	}

	return &Client{
		conn:    conn,
		channel: channel,
		queue:   queue,
	}, nil
}

// Close закрывает соединение с RabbitMQ сервером.
func (c *Client) Close() error {
	err := c.conn.Close()
	return errors.Wrap(err, "[rabbitmq::Close]: failed to close amqp connection")
}

// Publish публикует сообщение в очередь.
func (c *Client) Publish(ctx context.Context, msg []byte) error {
	err := c.channel.PublishWithContext(ctx,
		"",
		c.queue,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        msg,
		})

	return errors.Wrap(err, "[rabbitmq::Publish]: failed to publish amqp message")
}

// Consume запускает процесс непрерывного чтения опубликованных сообщений.
func (c *Client) Consume(ctx context.Context) (<-chan amqp.Delivery, error) {
	msgs, err := c.channel.ConsumeWithContext(ctx, c.queue, "", true, false, false, false, nil)
	return msgs, errors.Wrap(err, "[rabbitmq::Consume]: failed to receive messages from amqp")
}
