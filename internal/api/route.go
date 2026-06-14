package api

import (
	"github.com/gofiber/fiber/v2"
)

func (s *Server) Route(route fiber.Router) {
	route.Get("/ready/", s.DoPing)
	route.Post("/trip/", s.CreateTrip)
	route.Put("/trip/", s.MoveTripDraftToPublish)
	route.Get("/trip/:tripId", s.GetTrip)
	route.Get("/trip/events/:tripId", s.GetTripEvents)
}
