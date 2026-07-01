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
	app         *fiber.App // Поле должно быть здесь
	Registry    *prometheus.Registry
	Repository  *repository.RepoPg
	TripService *service.TripService
	ClientID    string
}

func NewServer(app *fiber.App,
	registry *prometheus.Registry,
	repo *repository.RepoPg,
	service *service.TripService,
	keycloakClientID string,
	isTest bool) *Server {
	s := &Server{
		app:         app,
		Registry:    registry,
		Repository:  repo,
		TripService: service,
		ClientID:    keycloakClientID,
	}
	if !isTest {
		s.app.Use(middleware.KeycloakRefreshTokenMiddleware(
			middleware.KeycloakConfig{
				Issuer:   configs.Env("KEYCLOAK_ISSUER", "http://localhost:8087/realms/sharetrip"),
				ClientID: keycloakClientID,
				// указываем правильный "secret"
				ClientSecret: configs.Env("KEYCLOAK_CLIENT_SECRET", "2Ac2dzdJ70w0vHiJYUX9RQ4LJxF8dhYV"),
			},
		))
	} else {
		s.app.Use(func(c *fiber.Ctx) error {
			// для теста имитируем наличие валидных Claims
			// для этого создаем структуру, которую ожидает ClaimsFromContext
			mockClaims := &middleware.KeycloakClaims{
				Subject: "test-user",
				ResourceAccess: map[string]struct {
					Roles []string `json:"roles"`
				}{
					keycloakClientID: {
						// роль должна совпадает с той, что проверяет RequireClientRole
						Roles: []string{"client"},
					},
				},
			}

			// Кладем именно под тем ключом, который прописан в константе KeycloakClaimsKey
			c.Locals(middleware.KeycloakClaimsKey, mockClaims)

			c.Locals("client_id", keycloakClientID)
			return c.Next()
		})
	}

	return s
}
