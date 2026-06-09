package service

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/google/uuid"
	"job4j.ru/share-trip/internal/domain"
	"job4j.ru/share-trip/internal/dto"
	"job4j.ru/share-trip/internal/repository"
	"strconv"
	"strings"
	"time"
)

type TripRequest dto.TripRequest

type CreateTripRequest dto.CreateTripRequest

type CreateTripResponse dto.CreateTripResponse

func CreateTripCommand(r *repository.RepoPg, c *fiber.Ctx, req CreateTripRequest) (*TripRequest, error) {
	id := uuid.New().String()
	now := time.Now()
	templateDate := "2006-01-02 15:04:05"
	departureTime, err := time.Parse(templateDate, req.DepartureTime)
	fmt.Println("Ошибка во время разбора времени:", err)
	fmt.Println("Разобранное время:", departureTime)
	if now.After(departureTime) {
		log.Errorw("CreateTripCommand", domain.ErrExpiredDate)
		return nil, domain.ErrExpiredDate
	}
	availableSeats, err := strconv.Atoi(req.AvailableSeats)
	if err != nil {
		log.Errorw("CreateTripCommand", domain.ErrAvailableSeatsNotInteger)
		return nil, domain.ErrAvailableSeatsNotInteger
	}
	if availableSeats < 1 {
		log.Errorw("CreateTripCommand", domain.ErrAvailableSeatsInvalidValue)
		return nil, domain.ErrAvailableSeatsInvalidValue
	}
	if strings.EqualFold(strings.Trim(req.FromPoint, " "), strings.Trim(req.ToPoint, " ")) {
		log.Errorw("CreateTripCommand", domain.ErrAvailableSeatsFromPointToPointEqual)
		return nil, domain.ErrAvailableSeatsFromPointToPointEqual
	}

	err = r.Create(c.Context(), dto.Trip{
		ID:             id,
		DriverId:       req.DriverId,
		FromPoint:      req.FromPoint,
		ToPoint:        req.ToPoint,
		DepartureTime:  req.DepartureTime,
		AvailableSeats: req.AvailableSeats,
	})
	if err != nil {
		log.Errorw("s.Repository.Create", err)
		return nil, err
	}

	res := TripRequest{
		ID:             id,
		DriverId:       req.DriverId,
		FromPoint:      req.FromPoint,
		ToPoint:        req.ToPoint,
		DepartureTime:  req.DepartureTime,
		AvailableSeats: req.AvailableSeats,
	}

	return &res, nil
}
