package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

type TripEventRequest struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type GetTripEventsResponse struct {
	TripEvents []TripEventRequest `json:"tripEvents"`
}

func (s *Server) GetTripEvents(c *fiber.Ctx) error {
	tripId := c.Params("tripId")
	if tripId == "" {
		return fiber.NewError(fiber.StatusBadRequest, "tripId is required")
	}

	items, err := s.Repository.EventList(c.Context(), tripId)
	if err != nil {
		log.Errorw("s.Repository.List", err)
		return fiber.NewError(fiber.StatusInternalServerError, "internal server error")
	}

	res := make([]TripEventRequest, 0, len(items))
	for _, item := range items {
		res = append(res, TripEventRequest{
			ID:   item.ID,
			Name: item.Name,
		})
	}

	return c.Status(fiber.StatusOK).JSON(GetTripEventsResponse{TripEvents: res})
}
