package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/google/uuid"
	tracker "job4j.ru/share-trip/internal/domain"
)

type ItemRequest struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type CreateItemRequest struct {
	Name string `json:"name"`
}

type CreateItemResponse struct {
	Item ItemRequest `json:"item"`
}

func (s *Server) CreateItem(c *fiber.Ctx) error {
	var req CreateItemRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid JSON body")
	}
	if req.Name == "" {
		return fiber.NewError(fiber.StatusBadRequest, "name is required")
	}

	id := uuid.New().String()
	err := s.Repository.Create(c.Context(), tracker.Item{
		ID:   id,
		Name: req.Name,
	})
	if err != nil {
		log.Errorw("s.Repository.Create", err)
		return fiber.NewError(fiber.StatusInternalServerError, "internal server error")
	}

	res := ItemRequest{
		ID:   id,
		Name: req.Name,
	}
	return c.Status(fiber.StatusCreated).JSON(CreateItemResponse{Item: res})
}
