package dto

type CreateTripRequest struct {
	DriverId       string `json:"driverId"`
	FromPoint      string `json:"fromPoint"`
	ToPoint        string `json:"toPoint"`
	DepartureTime  string `json:"departureTime"`
	AvailableSeats string `json:"availableSeats"`
}

type CreateTripResponse struct {
	Trip TripRequest `json:"trip"`
}
