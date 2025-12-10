package rabbitmq

import (
	"context"
	"fmt"
	"log"

	"github.com/bytedance/sonic"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/ride4Low/contracts/events"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
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

	tr := otel.Tracer("rabbitmq")
	ctx, span := tr.Start(ctx, "rabbitmq.publish",
		trace.WithAttributes(
			// attribute.String("messaging.system", "rabbitmq"),
			// attribute.String("messaging.destination", routingKey),
			attribute.String("messaging.routing_key", routingKey),
		),
	)
	defer span.End()

	jsonMsg, err := sonic.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %v", err)
	}

	headers := make(amqp.Table)
	otel.GetTextMapPropagator().Inject(ctx, AMQPHeadersCarrier(headers))

	msg := amqp.Publishing{
		Headers:      headers,
		ContentType:  "application/json",
		Body:         jsonMsg,
		DeliveryMode: amqp.Persistent,
	}

	if err = p.rmq.publish(ctx, TripExchange, routingKey, msg); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return fmt.Errorf("failed to publish message: %v", err)
	}
	return nil
}

// func (p *Publisher) Publish(ctx context.Context, exchange, routingKey string, body []byte) error {
// 	headers := make(amqp.Table)
// 	otel.GetTextMapPropagator().Inject(ctx, AMQPHeadersCarrier(headers))

// 	msg := amqp.Publishing{
// 		Headers:      headers,
// 		ContentType:  "application/json",
// 		Body:         body,
// 		DeliveryMode: amqp.Persistent,
// 	}
// 	return p.rmq.publish(ctx, exchange, routingKey, msg)
// }

func (r *RabbitMQ) publish(ctx context.Context, exchange, routingKey string, msg amqp.Publishing) error {
	return r.Channel.PublishWithContext(ctx,
		exchange,   // exchange
		routingKey, // routing key
		false,      // mandatory
		false,      // immediate
		msg,
	)
}
