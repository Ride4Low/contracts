package rabbitmq

import (
	"context"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

// Consumer represents a RabbitMQ consumer
type Consumer struct {
	rmq     *RabbitMQ
	handler MessageHandler
}

// MessageHandler defines the interface for handling messages
type MessageHandler interface {
	// Handle(ctx context.Context, body []byte) error
	Handle(context.Context, amqp.Delivery) error
}

// NewConsumer creates a new RabbitMQ consumer
func NewConsumer(rmq *RabbitMQ, handler MessageHandler) *Consumer {
	return &Consumer{
		rmq:     rmq,
		handler: handler,
	}
}

// Consume starts consuming messages from a queue
func (c *Consumer) Consume(ctx context.Context, queueName string) error {
	// Set prefetch count to 1 for fair dispatch
	// This tells RabbitMQ not to give more than one message to a service at a time.
	// The worker will only get the next message after it has acknowledged the previous one.
	err := c.rmq.Channel.Qos(
		1,     // prefetchCount: Limit to 1 unacknowledged message per consumer
		0,     // prefetchSize: No specific limit on message size
		false, // global: Apply prefetchCount to each consumer individually
	)
	if err != nil {
		return fmt.Errorf("failed to set QoS: %v", err)
	}

	msgs, err := c.rmq.Channel.Consume(
		queueName, // queue
		"",        // consumer
		false,     // auto-ack
		false,     // exclusive
		false,     // no-local
		false,     // no-wait
		nil,       // args
	)
	if err != nil {
		return fmt.Errorf("failed to consume messages: %v", err)
	}

	for {
		select {
		case <-ctx.Done():
			return nil
		case msg := <-msgs:
			if err := c.handler.Handle(ctx, msg); err != nil {
				fmt.Printf("Failed to handle message: %v\n", err)
				msg.Nack(false, false) // Don't requeue the message
			} else {
				msg.Ack(false)
			}
		}
	}

}
