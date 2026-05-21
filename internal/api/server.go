package api

import (
	"job4j.ru/share-trip/internal/storage"
)

type Server struct {
	Repository *storage.RepoPg
}

func NewServer(repo *storage.RepoPg) *Server {
	return &Server{Repository: repo}
}
