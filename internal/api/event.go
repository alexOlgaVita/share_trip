package api

type EventName string

const (
	EventPublished EventName = "trip_published"
	EventStarted   EventName = "trip_started"
)
