package rabbitmq

import (
	"context"
	"fmt"
	"log"

	"github.com/bytedance/sonic"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/ride4Low/contracts/events"
)

type Publisher struct {
	rmq *RabbitMQ
}

func NewPublisher(rmq *RabbitMQ) *Publisher {
	return &Publisher{
		rmq: rmq,
	}
}

func (p *Publisher) PublishMessage(ctx context.Context, routingKey string, message events.AmqpMessage) error {
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

	return p.rmq.publish(ctx, TripExchange, routingKey, msg)
}

func (p *Publisher) Publish(ctx context.Context, exchange, routingKey string, body []byte) error {
	msg := amqp.Publishing{
		ContentType:  "application/json",
		Body:         body,
		DeliveryMode: amqp.Persistent,
	}
	return p.rmq.publish(ctx, exchange, routingKey, msg)
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
