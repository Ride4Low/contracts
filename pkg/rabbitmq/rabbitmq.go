package rabbitmq

import (
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/ride4Low/contracts/events"
)

// Exchange names
const (
	TripExchange = "trip"
)

type RabbitMQ struct {
	conn    *amqp.Connection
	Channel *amqp.Channel
}

func NewRabbitMQ(uri string) (*RabbitMQ, error) {
	conn, err := amqp.Dial(uri)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %v", err)
	}

	channel, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("failed to create channel: %v", err)
	}

	rmq := &RabbitMQ{
		conn:    conn,
		Channel: channel,
	}

	if err := rmq.setupExchangesAndQueues(); err != nil {
		return nil, fmt.Errorf("failed to setup exchanges and queues: %v", err)
	}

	return rmq, nil

}

func (r *RabbitMQ) Close() error {
	if err := r.Channel.Close(); err != nil {
		return fmt.Errorf("failed to close channel: %v", err)
	}

	if err := r.conn.Close(); err != nil {
		return fmt.Errorf("failed to close connection: %v", err)
	}

	return nil
}

func (r *RabbitMQ) setupExchangesAndQueues() error {
	if err := r.Channel.ExchangeDeclare(
		TripExchange,
		amqp.ExchangeTopic,
		true,
		false,
		false,
		false,
		nil,
	); err != nil {
		return fmt.Errorf("failed to declare exchange: %v", err)
	}

	if err := r.declareAndBindQueue(
		events.FindAvailableDriversQueue,
		TripExchange,
		[]string{events.TripEventCreated},
	); err != nil {
		return fmt.Errorf("failed to declare and bind queue: %v", err)
	}

	return nil
}

func (r *RabbitMQ) declareAndBindQueue(queueName string, exchangeName string, routingKey []string) error {
	q, err := r.Channel.QueueDeclare(
		queueName, // name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments with DLX config
	)
	if err != nil {
		return fmt.Errorf("failed to declare queue: %v", err)
	}

	for _, key := range routingKey {
		if err := r.Channel.QueueBind(
			q.Name,
			key,
			exchangeName,
			false,
			nil,
		); err != nil {
			return fmt.Errorf("failed to bind queue: %v", err)
		}
	}

	return nil
}
