package rabbitmq

import (
	"context"
	"fmt"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
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

	go func() {
		for {
			select {
			case <-ctx.Done():
				log.Println("Consumer stopped")
				return
			case msg, ok := <-msgs:
				if !ok {
					log.Println("Channel closed")
					return
				}

				// Extract the context from the headers
				ctx := otel.GetTextMapPropagator().Extract(ctx, AMQPHeadersCarrier(msg.Headers))
				tr := otel.Tracer("rabbitmq")
				ctx, span := tr.Start(ctx, "rabbitmq.consume",
					trace.WithAttributes(
						// attribute.String("messaging.system", "rabbitmq"),
						attribute.String("messaging.destination", queueName),
						attribute.String("messaging.routing_key", msg.RoutingKey),
					),
				)

				if err := c.handler.Handle(ctx, msg); err != nil {
					fmt.Printf("Failed to handle message: %v\n", err)
					span.RecordError(err)
					msg.Nack(false, false) // Don't requeue the message
				} else {
					msg.Ack(false)
				}
				span.End()
			}
		}
	}()

	return nil

}
