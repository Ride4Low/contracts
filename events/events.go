package events

import "github.com/ride4Low/contracts/proto/trip"

// Event type constants
const (
	TripEventCreated = "trip.event.created"
)

// TripEventData is the payload for trip-related events
type TripEventData struct {
	Trip *trip.Trip `json:"trip"`
}
