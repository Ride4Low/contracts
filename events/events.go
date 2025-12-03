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
	FindAvailableDriversQueue       = "find_available_drivers"
	NotifyDriverNoDriversFoundQueue = "notify_driver_no_drivers_found"
	DriverCmdTripRequestQueue       = "driver_cmd_trip_request"
)

// Event type constants
const (
	// Trip events (trip.event.*)
	TripEventCreated        = "trip.event.created"
	TripEventNoDriversFound = "trip.event.no_drivers_found"

	// Driver commands (driver.cmd.*)
	DriverCmdTripRequest = "driver.cmd.trip_request"
)

// TripEventData is the payload for trip-related events
type TripEventData struct {
	Trip *trip.Trip `json:"trip"`
}
