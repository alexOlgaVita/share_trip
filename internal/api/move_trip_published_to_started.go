package api

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"job4j.ru/share-trip/internal/api/dto"
	"job4j.ru/share-trip/internal/domain"
	"strings"
)

type MoveTripPublishToStartedModel dto.MoveTripPublishToStartedModelRequest

type MoveTripPublishToStartedModelResponse dto.MoveTripPublishToStartedModelResponse

func (s *Server) MoveTripPublishedToStarted(c *fiber.Ctx) error {
	tracer := otel.Tracer("trip-api")
	ctx, span := tracer.Start(c.UserContext(), "MoveTripPublishToStartedTripHandler")
	defer span.End()

	c.Set("trace-id", span.SpanContext().TraceID().String())

	var req MoveTripPublishToStartedModel
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid JSON body")
	}

	if err := startedValidate(&req); err != nil {
		return err
	}

	span.SetAttributes(
		attribute.String("trip_id", req.TripID),
		attribute.String("driver_id", req.DriverId),
	)

	var resp, err = s.TripService.MoveTripPublishToStarted(ctx, dto.UpdateTripRequest{
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
		if errors.Is(err, domain.ErrStatusIsStartedAlready) {
			return fiber.NewError(fiber.StatusNoContent, "trip's status is started already")
		}
		if strings.Contains(err.Error(), domain.ErrNotAllowedToStart.Error()) {
			index := strings.Index(err.Error(), domain.ErrNotAllowedToStart.Error())
			if index != -1 {
				return fiber.NewError(fiber.StatusConflict, err.Error()[index:])
			} else {
				return fiber.NewError(fiber.StatusConflict, err.Error())
			}
		}
		log.Errorw("s.TripService.MovePublishToStart", err)
		return fiber.NewError(fiber.StatusInternalServerError, "internal server error")
	}

	return c.Status(fiber.StatusOK).JSON(resp)
}

func startedValidate(req *MoveTripPublishToStartedModel) error {
	if req.TripID == "" {
		return fiber.NewError(fiber.StatusBadRequest, "tripID is required")
	}
	if req.ClientID == "" {
		return fiber.NewError(fiber.StatusBadRequest, "clientID is required")
	}
	return nil
}
