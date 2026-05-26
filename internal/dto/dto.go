package dto

type TripRequest struct {
	ID             string `json:"id"`
	DriverId       string `json:"driverId"`
	FromPoint      string `json:"fromPoint"`
	ToPoint        string `json:"toPoint"`
	DepartureTime  string `json:"departureTime"`
	AvailableSeats string `json:"availableSeats"`
}

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
