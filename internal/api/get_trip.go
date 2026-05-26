package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

type GetTripResponse struct {
	Trip TripRequest `json:"trip"`
}

func (s *Server) GetTrip(c *fiber.Ctx) error {
	tripId := c.Params("tripId")
	if tripId == "" {
		return fiber.NewError(fiber.StatusBadRequest, "tripId is required")
	}

	trip, err := s.Repository.Get(c.Context(), tripId)
	if err != nil {
		log.Errorw("s.Repository.Get", err)
		return fiber.NewError(fiber.StatusInternalServerError, "internal server error")
	}

	res := TripRequest{
		ID:             trip.ID,
		DriverId:       trip.DriverId,
		FromPoint:      trip.FromPoint,
		ToPoint:        trip.ToPoint,
		DepartureTime:  trip.DepartureTime,
		AvailableSeats: trip.AvailableSeats,
	}

	return c.Status(fiber.StatusOK).JSON(GetTripResponse{Trip: res})
}
