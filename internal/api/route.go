package api

import "github.com/gofiber/fiber/v2"

func (s *Server) Route(route fiber.Router) {
	route.Get("/ready/", s.DoPing)
	route.Get("/count/", s.GetCount)
	route.Post("/item/", s.CreateItem)

}
