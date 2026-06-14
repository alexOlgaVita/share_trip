package api

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"job4j.ru/share-trip/internal/domain"
	"job4j.ru/share-trip/internal/dto"
)

type MoveTripDraftToPublishModelRequest dto.MoveTripDraftToPublishModelRequest

type MoveTripDraftToPublishModelResponse dto.MoveTripDraftToPublishModelResponse

func (s *Server) MoveTripDraftToPublish(c *fiber.Ctx) error {
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

	var resp, err = s.TripService.MoveTripDraftToPublish(c.Context(), dto.UpdateTripRequest{
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
