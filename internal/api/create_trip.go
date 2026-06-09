package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"job4j.ru/share-trip/internal/dto"
	"job4j.ru/share-trip/internal/service"
)

type TripRequest dto.TripRequest

type CreateTripRequest dto.CreateTripRequest

type CreateTripResponse dto.CreateTripResponse

func (s *Server) CreateTrip(c *fiber.Ctx) error {
	var req CreateTripRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid JSON body")
	}

	if req.DriverId == "" {
		return fiber.NewError(fiber.StatusBadRequest, "driverId is required")
	}
	if req.FromPoint == "" {
		return fiber.NewError(fiber.StatusBadRequest, "fromPoint is required")
	}
	if req.ToPoint == "" {
		return fiber.NewError(fiber.StatusBadRequest, "toPoint is required")
	}
	if req.DepartureTime == "" {
		return fiber.NewError(fiber.StatusBadRequest, "departureTime is required")
	}
	if req.AvailableSeats == "" {
		return fiber.NewError(fiber.StatusBadRequest, "availableSeats is required")
	}

	res, error := service.CreateTripCommand(s.Repository, c, service.CreateTripRequest(req))
	if error != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "internal server error")
	}
	return c.Status(fiber.StatusCreated).JSON(CreateTripResponse{Trip: dto.TripRequest(*res)})
}

func (s *Server) CreateTripNew(c *fiber.Ctx) error {
	var req CreateTripRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid JSON body")
	}

	if req.DriverId == "" {
		return fiber.NewError(fiber.StatusBadRequest, "driverId is required")
	}
	if req.FromPoint == "" {
		return fiber.NewError(fiber.StatusBadRequest, "fromPoint is required")
	}
	if req.ToPoint == "" {
		return fiber.NewError(fiber.StatusBadRequest, "toPoint is required")
	}
	if req.DepartureTime == "" {
		return fiber.NewError(fiber.StatusBadRequest, "departureTime is required")
	}
	if req.AvailableSeats == "" {
		return fiber.NewError(fiber.StatusBadRequest, "availableSeats is required")
	}

	var resp, err = s.TripService.CreateTrip(c.Context(), service.CreateTripRequest{
		DriverId:       req.DriverId,
		FromPoint:      req.FromPoint,
		ToPoint:        req.ToPoint,
		DepartureTime:  req.DepartureTime,
		AvailableSeats: req.AvailableSeats,
	})
	if err != nil {
		log.Errorw("s.TripService.CreateTrip", err)
		return fiber.NewError(fiber.StatusInternalServerError, "internal server error")
	}

	return c.Status(fiber.StatusCreated).JSON(resp)
}
