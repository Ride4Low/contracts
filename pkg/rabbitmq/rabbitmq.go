package rabbitmq

import (
	"context"
	"fmt"
	"log"

	"github.com/bytedance/sonic"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/ride4Low/contracts/events"
	"github.com/ride4Low/contracts/messaging"
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
		messaging.TripExchange,
		amqp.ExchangeTopic,
		true,
		false,
		false,
		false,
		nil,
	); err != nil {
		return fmt.Errorf("failed to declare exchange: %v", err)
	}

	r.declareAndBindQueue(
		messaging.FindAvailableDriversQueue,
		messaging.TripExchange,
		[]string{events.TripEventCreated},
	)

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
			exchangeName,
			key,
			false,
			nil,
		); err != nil {
			return fmt.Errorf("failed to bind queue: %v", err)
		}
	}

	return nil
}

func (r *RabbitMQ) PublishMessage(ctx context.Context, routingKey string, message messaging.AmqpMessage) error {
	log.Printf("Publishing message with routing key: %s", routingKey)

	jsonMsg, err := sonic.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %v", err)
	}

	msg := amqp.Publishing{
		ContentType:  "application/json",
		Body:         jsonMsg,
		DeliveryMode: amqp.Persistent,
	}

	return r.publish(ctx, messaging.TripExchange, routingKey, msg)
}

func (r *RabbitMQ) publish(ctx context.Context, exchange, routingKey string, msg amqp.Publishing) error {
	return r.Channel.PublishWithContext(ctx,
		exchange,   // exchange
		routingKey, // routing key
		false,      // mandatory
		false,      // immediate
		msg,
	)
}
