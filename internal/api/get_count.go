package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

type GetItemResponse struct {
	Count string `json:"count"`
}

func (s *Server) GetCount(c *fiber.Ctx) error {
	countRes, err := s.Repository.GetCount(c.Context())
	if err != nil {
		log.Errorw("s.Repository.GetCount", err)
		return fiber.NewError(fiber.StatusInternalServerError, "internal server error")
	}

	res := countRes

	return c.Status(fiber.StatusOK).JSON(GetItemResponse{Count: res})
}
