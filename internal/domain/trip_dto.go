package domain

type Trip struct {
	ID             string
	DriverID       string
	FromPoint      string
	ToPoint        string
	DepartureTime  string
	AvailableSeats string
	Status         string
	CreatedAt      string
}

type GetTripResponse struct {
	ID             string
	DriverID       string
	FromPoint      string
	ToPoint        string
	DepartureTime  string
	AvailableSeats string
	Status         string
	CreatedAt      string
}

type CreateTripRequest struct {
	ID             string
	DriverID       string
	FromPoint      string
	ToPoint        string
	DepartureTime  string
	AvailableSeats string
	Status         TripStatus
	CreatedAt      string
}

type CreateTripResponse struct {
	ID             string
	DriverID       string
	FromPoint      string
	ToPoint        string
	DepartureTime  string
	AvailableSeats string
	Status         TripStatus
	CreatedAt      string
}

type MoveTripDraftToPublishRequest struct {
	TripID   string
	ClientID string
}

type MoveTripDraftToPublishResponse struct {
	ID             string
	DriverID       string
	FromPoint      string
	ToPoint        string
	DepartureTime  string
	AvailableSeats string
	Status         string
	CreatedAt      string
}

type MoveTripPublishToStartRequest struct {
	TripID   string
	ClientID string
}

type MoveTripPublishToStartResponse struct {
	ID             string
	DriverID       string
	FromPoint      string
	ToPoint        string
	DepartureTime  string
	AvailableSeats string
	Status         string
	CreatedAt      string
}
