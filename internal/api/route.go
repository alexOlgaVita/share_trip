package api

import (
	"github.com/gofiber/adaptor/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"job4j.ru/share-trip/internal/middleware"
)

func (s *Server) Route(route fiber.Router) {
	route.Get("/ready/", s.DoPing)

	route.Post(
		"/trip/",
		middleware.RequireClientRole(s.ClientID, "client"),
		s.CreateTrip,
	)

	route.Put(
		"/trip/",
		middleware.RequireClientRole(s.ClientID, "client"),
		s.MoveTripDraftToPublish,
	)

	route.Get(
		"/trip/:tripId",
		middleware.RequireClientRole(s.ClientID, "client"),
		s.GetTrip,
	)

	route.Get("/trip/events/:tripId", s.GetTripEvents)
	route.Get("/metrics", adaptor.HTTPHandler(promhttp.HandlerFor(s.Registry, promhttp.HandlerOpts{})))
}
