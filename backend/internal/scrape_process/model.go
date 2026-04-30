package scrapeprocess

// Define the struct we want Claude to populate
type AIExtraction struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	ImageURL    string `json:"image_url"`

	OutboundDepartingLocation *string `json:"outbound_departing_location,omitempty"`
	OutboundArrivingLocation  *string `json:"outbound_arriving_location,omitempty"`
	OutboundDepartingAt       *string `json:"outbound_departing_at,omitempty"`
	OutboundArrivingAt        *string `json:"outbound_arriving_at,omitempty"`
	InboundDepartingLocation  *string `json:"inbound_departing_location,omitempty"`
	InboundArrivingLocation   *string `json:"inbound_arriving_location,omitempty"`
	InboundDepartingAt        *string `json:"inbound_departing_at,omitempty"`
	InboundArrivingAt         *string `json:"inbound_arriving_at,omitempty"`

	StartAt  *string `json:"start_at,omitempty"`
	EndAt    *string `json:"end_at,omitempty"`
	Location *string `json:"location,omitempty"`

	OutboundDurationMinutes *string `json:"outbound_duration_minutes,omitempty"`
	InboundDurationMinutes  *string `json:"inbound_duration_minutes,omitempty"`
}
