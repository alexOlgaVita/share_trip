package main

import (
	"context"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/prometheus/client_golang/prometheus"
	"job4j.ru/share-trip/configs"
	"job4j.ru/share-trip/internal/api"
	"job4j.ru/share-trip/internal/appl"
	"job4j.ru/share-trip/internal/domain"
	"job4j.ru/share-trip/internal/observability/metrics"
	"job4j.ru/share-trip/internal/observability/middleware"
	"job4j.ru/share-trip/internal/repository"
	"job4j.ru/share-trip/internal/service"
	"log"
	"os"
)

func main() {
	app := fiber.New()

	logger, logFile, err := appl.NewLogger()

	if err != nil {
		panic(err)
	}

	defer func() {
		if err := logFile.Close(); err != nil {
			fmt.Fprintf(os.Stderr, "failed to close log file: %v\n", err)
		}
	}()

	app.Use(middleware.Correlation(logger))

	ctx := context.Background()

	cfg := repository.Config{
		Host:     configs.Env("DB_HOST", "localhost"),
		Port:     configs.EnvInt("DB_PORT", 6543),
		User:     configs.Env("DB_USER", "postgres"),
		Password: configs.Env("DB_PASSWORD", "password"),
		DBName:   configs.Env("DB_NAME", "sharetrip"),
		SSLMode:  configs.Env("DB_SSLMODE", "disable"),
	}

	pool, err := repository.NewPool(ctx, cfg.DSN())
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	registry := prometheus.NewRegistry()
	m := metrics.New(registry)
	repo := repository.NewRepoPg(m, pool)
	srv := service.NewTripService(logger, m, pool, &domain.TripUsecase{
		TripRepo: repo,
	})

	server := api.NewServer(registry, repo, srv)
	app.Use(api.NewHTTPMetricsMiddleware(m))

	//server.Route(app.Group("/api"))
	server.Route(app.Group("/"))

	err = app.Listen(":8080")
	if err != nil {
		log.Fatal(err)
	}
}
