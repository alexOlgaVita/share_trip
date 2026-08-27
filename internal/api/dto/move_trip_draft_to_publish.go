package dto

type MoveTripDraftToPublishModelRequest struct {
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

type MoveTripDraftToPublishModelResponse struct {
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
