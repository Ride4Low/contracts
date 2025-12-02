package events

import "github.com/ride4Low/contracts/proto/trip"

type AmqpMessage struct {
	OwnerID string `json:"ownerId"`
	Data    any    `json:"data"`
}

// queues
const (
	FindAvailableDriversQueue = "find_available_drivers"
)

// events
const (
	TripEventCreated = "trip.event.created"
)

type TripEventData struct {
	Trip *trip.Trip `json:"trip"`
}
