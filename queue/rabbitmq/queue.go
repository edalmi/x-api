package rabbitmq

import (
	"context"
	"time"

	"github.com/edalmi/x-api/queue"
	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQ struct {
	conn *amqp.Connection
}

func (q *RabbitMQ) Push(ctx context.Context, name string, msg queue.Message) error {
	ch, err := q.conn.Channel()
	if err != nil {
		return err
	}

	defer ch.Close()

	qd, err := ch.QueueDeclare(
		name,
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = ch.PublishWithContext(ctx,
		"",
		qd.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: msg.ContentType,
			Body:        msg.Body,
		},
	)
	if err != nil {
		return err
	}

	return nil
}

func (q *RabbitMQ) Pull(ctx context.Context, name string) (<-chan *queue.Message, error) {
	ch, err := q.conn.Channel()
	if err != nil {
		return nil, err
	}

	defer ch.Close()

	qd, err := ch.QueueDeclare(
		name,
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, err
	}

	msgs, err := ch.Consume(
		qd.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, err
	}

	c := make(chan *queue.Message, 1)
	go func() {
		for d := range msgs {
			c <- &queue.Message{
				ContentType: d.ContentType,
				Body:        d.Body,
				AckFunc: func() error {
					orig := d
					return orig.Ack(false)
				},
			}
		}
	}()

	return c, nil
}

func (q *RabbitMQ) Close() error {
	return q.Close()
}
