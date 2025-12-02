package events

type AmqpMessage struct {
	OwnerID string `json:"ownerId"`
	Data    []byte `json:"data"`
}

// queues
const (
	FindAvailableDriversQueue = "find_available_drivers"
)

// events
const (
	TripEventCreated = "trip.event.created"
)
