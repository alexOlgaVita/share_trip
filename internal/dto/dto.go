package dto

type TripRequest struct {
	ID             string `json:"id"`
	DriverId       string `json:"driverId"`
	FromPoint      string `json:"fromPoint"`
	ToPoint        string `json:"toPoint"`
	DepartureTime  string `json:"departureTime"`
	AvailableSeats string `json:"availableSeats"`
	Status         string `json:"status"`
}

type Trip struct {
	ID             string
	DriverId       string
	FromPoint      string
	ToPoint        string
	DepartureTime  string
	AvailableSeats string
	Status         string
	CreatedAt      string
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

type SentNotificationTripPublishRequest struct {
	TripID string `json:"tripID"`
}

const (
	TripStatusDraft     = "draft"
	TripStatusPublished = "published"
	TripEventPublished  = "trip_published"
)

type TripEvent struct {
	ID   string
	Name string
}
