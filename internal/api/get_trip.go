package api

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"job4j.ru/share-trip/internal/domain"
)

type GetTripResponse struct {
	Trip TripRequest `json:"trip"`
}

func (s *Server) GetTrip(c *fiber.Ctx) error {
	tracer := otel.Tracer("trip-api")

	ctx, span := tracer.Start(c.UserContext(), "GetTripHandler")
	defer span.End()

	c.Set("trace-id", span.SpanContext().TraceID().String())

	tripId := c.Params("tripId")

	span.SetAttributes(
		attribute.String("trip_id", tripId),
	)

	if tripId == "" {
		return fiber.NewError(fiber.StatusBadRequest, "tripId is required")
	}

	trip, err := s.TripService.GetTrip(ctx, tripId)

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
