# Переменные
GO := go
GO_PKG := ./...
COVERAGE_PROFILE := coverage.out
GOOSE_MIGRATION_DIR := migrations
NAME_TRIP := create_table_trip
NAME_OUTBOX_EVENT := create_table_outbox_event
DOCKER_NAME := pg
BINARY_NAME := sharetrip.exe
MAIN_PKG := ./cmd/sharetrip
BIN_DIR := bin

# Цель по умолчанию
.PHONY: help
help:
	@echo "Available targets:"
	@echo "  test        - Run all *_test.go with module coverage (same as stand)"
	@echo "  coverage    - Run tests, print total, generate HTML coverage report"
	@echo "  cover       - Alias for coverage"
	@echo "  lint        - Run golangci-lint"
	@echo "  all         - Run lint, tests and coverage"
	@echo "  help        - Show this help"

# Как на стенде: все *_test.go, покрытие по всему модулю, итог total
.PHONY: test
test:
	$(GO) test -v -coverpkg=$(GO_PKG) -coverprofile=$(COVERAGE_PROFILE) $(GO_PKG)
	$(GO) tool cover -func=$(COVERAGE_PROFILE)

# HTML-отчёт по уже собранному coverage.out (тесты гоняются в test)
.PHONY: coverage cover
coverage cover: test
	$(GO) tool cover -html=$(COVERAGE_PROFILE) -o coverage.html
	@echo Coverage report generated: coverage.html

.PHONY: cover-report
cover-report:
	$(GO) test -cover -coverpkg=$(GO_PKG) $(GO_PKG)

# Запуск всех тестов
.PHONY: fmt
fmt:
	$(GO) fmt $(GO_PKG)

# Проверка кода с помощью golangci-lint
.PHONY: lint
lint:
	golangci-lint run

.PHONY: build
build:
	$(GO) build -o $(BIN_DIR)/$(BINARY_NAME) $(MAIN_PKG)

.PHONY: run
run:
	$(GO) run ./$(MAIN_PKG)/main.go

.PHONY: deps
deps:
	$(GO) mod tidy
	$(GO) mod download

.PHONY: migrate-status
migrate-status:
	goose -dir $(GOOSE_MIGRATION_DIR) postgres "postgres://postgres:password@localhost:6543/sharetrip?sslmode=disable" status

.PHONY: migrate-up
migrate-up:
	goose -dir $(GOOSE_MIGRATION_DIR) postgres "postgres://postgres:password@localhost:6543/sharetrip?sslmode=disable" up

.PHONY: migrate-down
migrate-down:
	goose -dir $(GOOSE_MIGRATION_DIR) postgres "postgres://postgres:password@localhost:6543/sharetrip?sslmode=disable" down

.PHONY: migrate-create
migrate-create:
	goose -dir $(GOOSE_MIGRATION_DIR) create -s $(NAME_TRIP) sql
	goose -dir $(GOOSE_MIGRATION_DIR) create -s $(NAME_OUTBOX_EVENT) sql

.PHONY: up
up:
	docker-compose up -d

.PHONY: down
down:
	docker-compose down --volumes

# Запуск всех проверок (тесты один раз: цель coverage зависит от test)
.PHONY: check
check: fmt lint coverage

.PHONY: e2eCheck
e2eCheck:
	curl http://localhost:8080/api/ready