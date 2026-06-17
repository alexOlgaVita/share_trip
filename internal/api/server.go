package api

import (
	"github.com/prometheus/client_golang/prometheus"
	"job4j.ru/share-trip/internal/repository"
	"job4j.ru/share-trip/internal/service"
)

type Server struct {
	Registry    *prometheus.Registry
	Repository  *repository.RepoPg
	TripService *service.TripService
}

func NewServer(registry *prometheus.Registry, repo *repository.RepoPg, service *service.TripService) *Server {
	return &Server{
		Registry:    registry,
		Repository:  repo,
		TripService: service,
	}
}
