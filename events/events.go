package events

import (
	"encoding/json"

	"github.com/ride4Low/contracts/proto/trip"
)

// AmqpMessage is the standard message format for RabbitMQ messages
type AmqpMessage struct {
	OwnerID string          `json:"ownerId"`
	Data    json.RawMessage `json:"data"`
}

// Queue names
const (
	FindAvailableDriversQueue = "find_available_drivers"
)

// Event type constants
const (
	TripEventCreated = "trip.event.created"
)

// TripEventData is the payload for trip-related events
type TripEventData struct {
	Trip *trip.Trip `json:"trip"`
}
