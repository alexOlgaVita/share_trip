package main

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"job4j.ru/share-trip/configs"
	"job4j.ru/share-trip/internal/api"
	"job4j.ru/share-trip/internal/repository"
	"log"
)

func main() {
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

	repo := repository.NewRepoPg(pool)

	server := api.NewServer(repo)

	app := fiber.New()
	server.Route(app.Group("/api"))

	err = app.Listen(":8080")
	if err != nil {
		log.Fatal(err)
	}
}
