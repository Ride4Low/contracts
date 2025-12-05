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
	DriverTripResponseQueue         = "driver_trip_response"
)

// Event type constants
const (
	// Trip events (trip.event.*)
	TripEventCreated        = "trip.event.created"
	TripEventNoDriversFound = "trip.event.no_drivers_found"

	// Driver commands (driver.cmd.*)
	DriverCmdTripRequest = "driver.cmd.trip_request"
	DriverCmdTripAccept  = "driver.cmd.trip_accept"
	DriverCmdTripDecline = "driver.cmd.trip_decline"
	DriverCmdLocation    = "driver.cmd.location"
)

// TripEventData is the payload for trip-related events
type TripEventData struct {
	Trip *trip.Trip `json:"trip"`
}

// WS Events

const (
	// DriverWSRegister = "driver.ws.register"
	DriverCmdRegister = "driver.cmd.register"
)
