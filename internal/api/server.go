package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/prometheus/client_golang/prometheus"
	"job4j.ru/share-trip/configs"
	"job4j.ru/share-trip/internal/middleware"
	"job4j.ru/share-trip/internal/repository"
	"job4j.ru/share-trip/internal/service"
)

type Server struct {
	app          *fiber.App // Поле должно быть здесь
	Registry     *prometheus.Registry
	Repository   *repository.RepoPg
	TripService  *service.TripService
	ClientID     string
	requiredRole string
	ContractCfg  configs.Contract
}

func NewServer(app *fiber.App,
	registry *prometheus.Registry,
	repo *repository.RepoPg,
	service *service.TripService,
	kc configs.Keycloak,
	testOverride middleware.KeycloakConfig,
	contractCfg configs.Contract,
) *Server {
	s := &Server{
		app:          app,
		Registry:     registry,
		Repository:   repo,
		TripService:  service,
		ClientID:     kc.ClientID,
		requiredRole: kc.RequiredRole,
		ContractCfg:  contractCfg,
	}

	mwCfg := middleware.KeycloakConfig{
		Issuer:       kc.Issuer,
		ClientID:     kc.ClientID,
		ClientSecret: kc.ClientSecret,
	}

	if testOverride.HTTPClient != nil {
		mwCfg = testOverride
	}

	s.app.Use(middleware.KeycloakRefreshTokenMiddleware(mwCfg))
	return s
}
