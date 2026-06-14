package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"job4j.ru/share-trip/internal/dto"
	"strconv"
	"strings"
	"time"
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
	now := time.Now()
	templateDate := "2006-01-02 15:04:05"
	departureTime, err := time.Parse(templateDate, req.DepartureTime)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "date format error")
	}
	if now.After(departureTime) {
		return fiber.NewError(fiber.StatusBadRequest, "date is expired")
	}
	availableSeats, err := strconv.Atoi(req.AvailableSeats)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "attribute 'availableSeats' doesn't match to integer type")
	}
	if availableSeats < 1 {
		return fiber.NewError(fiber.StatusBadRequest, "attribute 'availableSeats' has value less 1")
	}
	if strings.EqualFold(strings.Trim(req.FromPoint, " "), strings.Trim(req.ToPoint, " ")) {
		return fiber.NewError(fiber.StatusBadRequest, "fromPoint and toPoint have the same value")
	}

	resp, err := s.TripService.CreateTrip(c.Context(), dto.CreateTripRequest{
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
