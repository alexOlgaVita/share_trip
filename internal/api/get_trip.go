package api

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"job4j.ru/share-trip/internal/domain"
)

type GetTripResponse struct {
	Trip TripRequest `json:"trip"`
}

func (s *Server) GetTrip(c *fiber.Ctx) error {
	tripId := c.Params("tripId")
	if tripId == "" {
		return fiber.NewError(fiber.StatusBadRequest, "tripId is required")
	}

	trip, err := s.TripService.GetTrip(c.Context(), tripId)

	if err != nil {
		if errors.Is(err, domain.ErrTripNotFound) {
			return fiber.NewError(fiber.StatusNotFound, "trip is not found")
		}
		log.Errorw(
			"get trip failed",
			"error", err,
			"trip_id", tripId,
		)
		return fiber.NewError(fiber.StatusInternalServerError, "internal server error")
	}

	res := TripRequest{
		ID:             trip.ID,
		DriverId:       trip.DriverId,
		FromPoint:      trip.FromPoint,
		ToPoint:        trip.ToPoint,
		DepartureTime:  trip.DepartureTime,
		AvailableSeats: trip.AvailableSeats,
		Status:         trip.Status,
	}

	return c.Status(fiber.StatusOK).JSON(GetTripResponse{Trip: res})
}
