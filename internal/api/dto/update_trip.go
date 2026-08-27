package dto

type UpdateTripRequest struct {
	DriverId       string `json:"driverId"`
	FromPoint      string `json:"fromPoint"`
	ToPoint        string `json:"toPoint"`
	DepartureTime  string `json:"departureTime"`
	AvailableSeats string `json:"availableSeats"`

	TripID   string
	ClientID string
	Status   string
}
