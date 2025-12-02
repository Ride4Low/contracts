package messaging

// AmqpMessage is the standard message format for RabbitMQ messages
type AmqpMessage struct {
	OwnerID string `json:"ownerId"`
	Data    any    `json:"data"`
}

// Exchange names
const (
	TripExchange = "trip"
)

// Queue names
const (
	FindAvailableDriversQueue = "find_available_drivers"
)
