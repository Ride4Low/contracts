package events

import (
	"encoding/json"

	"github.com/ride4Low/contracts/proto/driver"
	"github.com/ride4Low/contracts/proto/trip"
)

// AmqpMessage is the standard message format for RabbitMQ messages
type AmqpMessage struct {
	OwnerID string          `json:"ownerId"`
	Data    json.RawMessage `json:"data"`
}

// Queue names
const (
	FindAvailableDriversQueue        = "find_available_drivers"
	NotifyDriverNoDriversFoundQueue  = "notify_driver_no_drivers_found"
	DriverCmdTripRequestQueue        = "driver_cmd_trip_request"
	DriverTripResponseQueue          = "driver_trip_response"
	NotifyDriverAssignQueue          = "notify_driver_assign"
	PaymentTripResponseQueue         = "payment_trip_response"
	NotifyPaymentSessionCreatedQueue = "notify_payment_session_created"
	NotifyPaymentSuccessQueue        = "notify_payment_success"
	DeadLetterQueue                  = "dead_letter"
)

// Event type constants
const (
	// Trip events (trip.event.*)
	TripEventCreated             = "trip.event.created"
	TripEventNoDriversFound      = "trip.event.no_drivers_found"
	TripEventDriverNotInterested = "trip.event.driver_not_interested"
	TripEventDriverAssigned      = "trip.event.driver_assigned"

	// Driver commands (driver.cmd.*)
	DriverCmdTripRequest = "driver.cmd.trip_request"
	DriverCmdTripAccept  = "driver.cmd.trip_accept"
	DriverCmdTripDecline = "driver.cmd.trip_decline"
	DriverCmdLocation    = "driver.cmd.location"

	// Payment commands (payment.cmd.*)
	PaymentCmdCreateSession = "payment.cmd.create_session"
	PaymentCmdSelectCard    = "payment.cmd.select_card"
	PaymentCmdSelectCrypto  = "payment.cmd.select_crypto"

	// Payment events (payment.event.*)
	PaymentEventSessionCreated = "payment.event.session_created"
	PaymentEventSuccess        = "payment.event.success"
)

// TripEventData is the payload for trip-related events
type TripEventData struct {
	Trip *trip.Trip `json:"trip"`
}

type DriverTripResponseData struct {
	Driver  *driver.Driver `json:"driver"`
	TripID  string         `json:"tripID"`
	RiderID string         `json:"riderID"`
}

// WS Events
const (
	// DriverWSRegister = "driver.ws.register"
	DriverCmdRegister = "driver.cmd.register"
)

type PaymentTripResponseData struct {
	TripID   string  `json:"tripID"`
	UserID   string  `json:"userID"`
	DriverID string  `json:"driverID"`
	Amount   float64 `json:"amount"`
	Currency string  `json:"currency"`
}

type PaymentEventSessionCreatedData struct {
	TripID    string  `json:"tripID"`
	SessionID string  `json:"sessionID"`
	Amount    float64 `json:"amount"`
	Currency  string  `json:"currency"`
}

type PaymentStatusUpdateData struct {
	TripID   string `json:"tripID"`
	UserID   string `json:"userID"`
	DriverID string `json:"driverID"`
}

type PaymentSelectCardData struct {
	TripID string `json:"tripID"`
	UserID string `json:"userID"`
}
