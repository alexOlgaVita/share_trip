package service

type EventName string

const (
	TripPublished EventName = "trip_published"
	TripStarted   EventName = "trip_started"
)
