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
	"job4j.ru/share-trip/internal/observability/tracing"
	"job4j.ru/share-trip/internal/repository"
	"job4j.ru/share-trip/internal/service"
	"log"
	"os"
	"time"
)

func main() {
	ctx := context.Background()
	app := fiber.New()

	// Настройка логгера и middleware
	logger, logFile, err := appl.NewLogger()

	if err != nil {
		panic(err)
	}

	defer func() {
		if err := logFile.Close(); err != nil {
			fmt.Fprintf(os.Stderr, "failed to close log file: %v\n", err)
		}
	}()

	// Подключение Middleware для автоматического захвата трасс Fiber
	app.Use(middleware.Correlation(logger))

	// Инициализация трассировки
	tp, err := tracing.NewProvider(ctx, tracing.Config{
		ServiceName:    "share-trip",
		ServiceVersion: "1.0.0",
		Environment:    "local",
		Endpoint:       "localhost:4319",
	})
	if err != nil {
		logger.Error("init tracing failed", "error", err)
		os.Exit(1)
	}

	defer func() {
		shutdownCtx, cancel := context.WithTimeout(
			context.Background(),
			5*time.Second,
		)
		defer cancel()

		if err := tp.Shutdown(shutdownCtx); err != nil {
			logger.Error("shutdown tracing failed", "error", err)
		}
	}()

	// Инициализация репозитория
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

	// Инициализация сервиса прометеус
	registry := prometheus.NewRegistry()
	m := metrics.New(registry)
	repo := repository.NewRepoPg(m, pool)
	srv := service.NewTripService(logger, m, pool, &domain.TripUsecase{
		TripRepo: repo,
	})

	server := api.NewServer(registry, repo, srv)

	app.Use(api.NewHTTPMetricsMiddleware(m))

	// Настройка роута
	//server.Route(app.Group("/api"))
	server.Route(app.Group("/"))

	err = app.Listen(":8080")
	if err != nil {
		log.Fatal(err)
	}
}
