package repository

type TripEvent struct {
	ID   string
	Name string
}

type SentNotificationTripPublishRequest struct {
	TripID string `json:"tripID"`
}
