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

type SentNotificationTripPublishRequest struct {
	TripID string `json:"tripID"`
}

const (
	TripStatusDraft            = "draft"
	TripStatusPublished        = "published"
	TripStatusStarted          = "started"
	TripEventPublished         = "trip_published"
	TripEventStarted           = "trip_started"
	ReasonNoActiveContract     = "NO_ACTIVE_CONTRACT"
	ReasonContractSuspended    = "CONTRACT_SUSPENDED"
	ReasonContractTerminated   = "CONTRACT_TERMINATED"
	ReasonContractNotStarted   = "CONTRACT_NOT_STARTED"
	ReasonContractExpired      = "CONTRACT_EXPIRED"
	ReasonServiceNotInContract = "SERVICE_NOT_IN_CONTRACT"
	ReasonServiceDisabled      = "SERVICE_DISABLED"
)

type TripEvent struct {
	ID   string
	Name string
}
