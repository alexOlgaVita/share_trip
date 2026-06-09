package api

import (
	"job4j.ru/share-trip/internal/repository"
	"job4j.ru/share-trip/internal/service"
)

type Server struct {
	Repository  *repository.RepoPg
	TripService *service.TripService
}

func NewServer(repo *repository.RepoPg) *Server {
	return &Server{Repository: repo}
}
