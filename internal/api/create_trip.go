package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"job4j.ru/share-trip/internal/observability/logctx"
	"log/slog"
	"strconv"
	"strings"
	"time"
)

func (s *Server) CreateTrip(c *fiber.Ctx) error {
	tracer := otel.Tracer("trip-api")
	ctx, span := tracer.Start(c.UserContext(), "CreateTripHandler")
	defer span.End()

	c.Set("trace-id", span.SpanContext().TraceID().String())

	logger := logctx.Logger(ctx).With(
		slog.String("server", "TripServer"),
		slog.String("handler", "CreateTrip"),
	)

	var req CreateTripRequest
	if err := c.BodyParser(&req); err != nil {
		logger.Warn(
			"create trip failed: invalid json body",
			slog.Any("error", err),
		)
		return fiber.NewError(fiber.StatusBadRequest, "invalid JSON body")
	}

	span.SetAttributes(
		attribute.String("driver_id", req.DriverID),
		attribute.String("departure_time", req.DepartureTime),
		attribute.String("to_point", req.ToPoint),
		attribute.String("from_point", req.FromPoint),
		attribute.String("available_seats", req.AvailableSeats),
	)

	if err := createValidate(&req, logger); err != nil {
		return err
	}

	logger = logger.With(
		slog.String("driverId", req.DriverID),
		slog.String("fromPoint", req.FromPoint),
		slog.String("toPoint", req.ToPoint),
		slog.String("departureTime", req.DepartureTime),
		slog.String("availableSeats", req.AvailableSeats),
	)
	ctx = logctx.WithLogger(ctx, logger)
	logger.Info("create trip request accepted")

	resp, err := s.TripService.CreateTrip(ctx, toServiceCreateTripRequest(req))
	if err != nil {
		log.Errorw("s.TripService.CreateTrip", err)
		logger.Error(
			"create trip failed",
			slog.Any("error", err),
		)
		return fiber.NewError(fiber.StatusInternalServerError, "internal server error")
	}

	logger.Info(
		"create trip completed",
		slog.String("trip_id", resp.ID),
	)

	apiResp, err := fromServiceCreateTripResponse(resp)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "internal mapping error")
	}

	return c.Status(fiber.StatusCreated).JSON(apiResp)
}

func createValidate(req *CreateTripRequest, logger *slog.Logger) error {
	if req.DriverID == "" {
		logger.Warn("create trip failed: driverId is required")
		return fiber.NewError(fiber.StatusBadRequest, "driverId is required")
	}

	if req.FromPoint == "" {
		logger.Warn("create trip failed: fromPoint is required")
		return fiber.NewError(fiber.StatusBadRequest, "fromPoint is required")
	}

	if req.ToPoint == "" {
		logger.Warn("create trip failed: toPoint is required")
		return fiber.NewError(fiber.StatusBadRequest, "toPoint is required")
	}

	if req.DepartureTime == "" {
		logger.Warn("create trip failed: departureTime is required")
		return fiber.NewError(fiber.StatusBadRequest, "departureTime is required")
	}

	if req.AvailableSeats == "" {
		logger.Warn("create trip failed: availableSeats is required")
		return fiber.NewError(fiber.StatusBadRequest, "availableSeats is required")
	}

	now := time.Now()
	templateDate := "2006-01-02 15:04:05"
	departureTime, err := time.Parse(templateDate, req.DepartureTime)
	if err != nil {
		logger.Warn("departureTime has bad date format")
		return fiber.NewError(fiber.StatusBadRequest, "date format error")
	}

	if now.After(departureTime) {
		logger.Warn("departureTime is expired date")
		return fiber.NewError(fiber.StatusBadRequest, "departureTime is expired date")
	}

	availableSeats, err := strconv.Atoi(req.AvailableSeats)
	if err != nil {
		logger.Warn("attribute 'availableSeats' doesn't match to integer type")
		return fiber.NewError(fiber.StatusBadRequest, "attribute 'availableSeats' doesn't match to integer type")
	}

	if availableSeats < 1 {
		logger.Warn("attribute 'availableSeats' has value less 1")
		return fiber.NewError(fiber.StatusBadRequest, "attribute 'availableSeats' has value less 1")
	}

	if strings.EqualFold(strings.Trim(req.FromPoint, " "), strings.Trim(req.ToPoint, " ")) {
		logger.Warn("fromPoint and toPoint have the same value")
		return fiber.NewError(fiber.StatusBadRequest, "fromPoint and toPoint have the same value")
	}
	return nil
}
