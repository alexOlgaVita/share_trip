package api

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"job4j.ru/share-trip/internal/domain"
	"job4j.ru/share-trip/internal/dto"
)

type MoveTripDraftToPublishModelRequest dto.MoveTripDraftToPublishModelRequest

type MoveTripDraftToPublishModelResponse dto.MoveTripDraftToPublishModelResponse

func (s *Server) MoveTripDraftToPublish(c *fiber.Ctx) error {
	tracer := otel.Tracer("trip-api")
	ctx, span := tracer.Start(c.UserContext(), "MoveTripDraftToPublishTripHandler")
	defer span.End()

	c.Set("trace-id", span.SpanContext().TraceID().String())

	var req MoveTripDraftToPublishModelRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid JSON body")
	}

	if req.TripID == "" {
		return fiber.NewError(fiber.StatusBadRequest, "tripID is required")
	}
	if req.ClientID == "" {
		return fiber.NewError(fiber.StatusBadRequest, "clientID is required")
	}

	span.SetAttributes(
		attribute.String("trip_id", req.TripID),
		attribute.String("driver_id", req.DriverId),
	)

	var resp, err = s.TripService.MoveTripDraftToPublish(ctx, dto.UpdateTripRequest{
		TripID:   req.TripID,
		ClientID: req.ClientID,
	})

	if err != nil {

		if errors.Is(err, domain.ErrTripNotFound) {
			return fiber.NewError(fiber.StatusNotFound, "trip is not found")
		}
		if errors.Is(err, domain.ErrClientNotDriver) {
			return fiber.NewError(fiber.StatusForbidden, "client is not driver of this trip")
		}
		if errors.Is(err, domain.ErrNotAllowedCurrentStatusToPublish) {
			return fiber.NewError(fiber.StatusConflict, "current status is not allowed for publish")
		}
		if errors.Is(err, domain.ErrStatusIsPublishedAlready) {
			return fiber.NewError(fiber.StatusNoContent, "trip's status is published already")
		}

		log.Errorw("s.TripService.MoveTripDraftToPublish", err)
		return fiber.NewError(fiber.StatusInternalServerError, "internal server error")
	}

	return c.Status(fiber.StatusOK).JSON(resp)
}
