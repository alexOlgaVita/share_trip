package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
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
		log.Errorw("s.TripService.MoveTripDraftToPublish", err)
		return fiber.NewError(fiber.StatusInternalServerError, "internal server error")
	}

	if err == nil && resp == nil {
		log.Errorw("s.TripService.MoveTripDraftToPublish", err)
		return fiber.NewError(fiber.StatusInternalServerError, "internal server error^ empty result")
	}

	return c.Status(fiber.StatusOK).JSON(resp)
}
