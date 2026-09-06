package api

type TripStatus string

const (
	StatusDraft     TripStatus = "draft"
	StatusPublished TripStatus = "published"
	StatusStarted   TripStatus = "started"
	StatusActive    TripStatus = "active"
	StatusClosed    TripStatus = "closed"
)
