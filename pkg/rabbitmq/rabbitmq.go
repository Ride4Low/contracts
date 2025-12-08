package rabbitmq

import (
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/ride4Low/contracts/events"
)

// Exchange names
const (
	TripExchange       = "trip"
	DeadLetterExchange = "dlx"
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
	var errs []error
	if err := r.Channel.Close(); err != nil {
		errs = append(errs, fmt.Errorf("failed to close channel: %v", err))
	}

	if err := r.conn.Close(); err != nil {
		errs = append(errs, fmt.Errorf("failed to close connection: %v", err))
	}

	if len(errs) > 0 {
		return fmt.Errorf("errors closing rabbitmq: %v", errs)
	}

	return nil
}

func (r *RabbitMQ) setupExchangesAndQueues() error {
	var topology = map[string][]struct {
		queueName   string
		routingKeys []string
	}{
		DeadLetterExchange: {
			{
				queueName:   events.DeadLetterQueue,
				routingKeys: []string{"#"}, // wildcard routing key to catch all messages
			},
		},
		TripExchange: {
			{
				queueName:   events.FindAvailableDriversQueue,
				routingKeys: []string{events.TripEventCreated, events.TripEventDriverNotInterested},
			},
			{
				queueName:   events.NotifyDriverNoDriversFoundQueue,
				routingKeys: []string{events.TripEventNoDriversFound},
			},
			{
				queueName:   events.DriverCmdTripRequestQueue,
				routingKeys: []string{events.DriverCmdTripRequest},
			},
			{
				queueName:   events.DriverTripResponseQueue,
				routingKeys: []string{events.DriverCmdTripAccept, events.DriverCmdTripDecline},
			},
			{
				queueName:   events.NotifyDriverAssignQueue,
				routingKeys: []string{events.TripEventDriverAssigned},
			},
			{
				queueName:   events.PaymentTripResponseQueue,
				routingKeys: []string{events.PaymentCmdCreateSession},
			},
			{
				queueName:   events.NotifyPaymentSessionCreatedQueue,
				routingKeys: []string{events.PaymentEventSessionCreated},
			},
			{
				queueName:   events.NotifyPaymentSuccessQueue,
				routingKeys: []string{events.PaymentEventSuccess},
			},
		},
	}

	for exchange, queues := range topology {
		if err := r.Channel.ExchangeDeclare(
			exchange,
			amqp.ExchangeTopic,
			true,
			false,
			false,
			false,
			nil,
		); err != nil {
			return fmt.Errorf("failed to declare exchange: %v", err)
		}

		for _, queueAndRoutingKeys := range queues {
			if err := r.declareAndBindQueue(
				queueAndRoutingKeys.queueName,
				exchange,
				queueAndRoutingKeys.routingKeys,
			); err != nil {
				return fmt.Errorf("failed to declare and bind queue: %v", err)
			}
		}
	}

	return nil
}

func (r *RabbitMQ) declareAndBindQueue(queueName string, exchangeName string, routingKeys []string) error {
	// Add dead letter configuration
	var args amqp.Table

	if queueName != events.DeadLetterQueue {
		args = amqp.Table{
			"x-dead-letter-exchange": DeadLetterExchange,
		}
	} else {
		args = amqp.Table{
			"x-message-ttl": 86400000, // 1 day in milliseconds
			// "x-message-ttl": 60000, // 1 minute
		}
	}

	q, err := r.Channel.QueueDeclare(
		queueName, // name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		args,      // arguments with DLX config
	)
	if err != nil {
		return fmt.Errorf("failed to declare queue: %v", err)
	}

	for _, key := range routingKeys {
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
