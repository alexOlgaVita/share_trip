package repository

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
	Status         string
	CreatedAt      string
}

type CreateTripResponse struct {
	ID             string
	DriverID       string
	FromPoint      string
	ToPoint        string
	DepartureTime  string
	AvailableSeats string
	Status         string
	CreatedAt      string
}

type UpdateTripRequest struct {
	ID             string
	DriverID       string
	FromPoint      string
	ToPoint        string
	DepartureTime  string
	AvailableSeats string
	Status         string
	CreatedAt      string
}

type UpdateTripResponse struct {
	ID             string
	DriverID       string
	FromPoint      string
	ToPoint        string
	DepartureTime  string
	AvailableSeats string
	Status         string
	CreatedAt      string
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
	ID             string
	DriverID       string
	FromPoint      string
	ToPoint        string
	DepartureTime  string
	AvailableSeats string
	Status         string
	CreatedAt      string
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
