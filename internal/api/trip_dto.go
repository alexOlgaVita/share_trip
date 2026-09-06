package api

type TripRequest struct {
	ID             string `json:"id"`
	DriverID       string `json:"driverId"`
	FromPoint      string `json:"fromPoint"`
	ToPoint        string `json:"toPoint"`
	DepartureTime  string `json:"departureTime"`
	AvailableSeats string `json:"availableSeats"`
	Status         string `json:"status"`
	CreatedAt      string `json:"createdAt"` //?убрать
}

type TripResponse struct {
	ID             string `json:"id"`
	DriverID       string `json:"driverId"`
	FromPoint      string `json:"fromPoint"`
	ToPoint        string `json:"toPoint"`
	DepartureTime  string `json:"departureTime"`
	AvailableSeats string `json:"availableSeats"`
	Status         string `json:"status"`
	CreatedAt      string `json:"createdAt"`
}

type GetTripResponse struct {
	Trip TripRequest `json:"trip"`
}

type CreateTripRequest struct {
	DriverID       string `json:"driverId"`
	FromPoint      string `json:"fromPoint"`
	ToPoint        string `json:"toPoint"`
	DepartureTime  string `json:"departureTime"`
	AvailableSeats string `json:"availableSeats"`
}

type CreateTripResponse struct {
	Trip TripResponse `json:"trip"`
}

type MoveTripDraftToPublishRequest struct {
	TripID   string `json:"tripID"`
	ClientID string `json:"clientID"`
}

type MoveTripDraftToPublishResponse struct {
	Trip TripResponse `json:"trip"`
}

type MoveTripPublishToStartRequest struct {
	TripID   string `json:"tripID"`
	ClientID string `json:"clientID"`
}

type MoveTripPublishToStartResponse struct {
	Trip TripResponse `json:"trip"`
}
