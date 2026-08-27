package dto

type MoveTripPublishToStartedModelRequest struct {
	ID             string
	DriverId       string
	FromPoint      string
	ToPoint        string
	DepartureTime  string
	AvailableSeats string
	Status         string
	CreatedAt      string

	TripID   string
	ClientID string
}

type MoveTripPublishToStartedModelResponse struct {
	ID             string
	DriverId       string
	FromPoint      string
	ToPoint        string
	DepartureTime  string
	AvailableSeats string
	Status         string
	CreatedAt      string

	TripID   string
	ClientID string
}
