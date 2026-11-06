package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"
	events "job4j.ru/share-trip/internal/kafka"
)

// Publisher мокирует events.Publisher.
type Publisher struct {
	mock.Mock
}

func (m *Publisher) PublishTripPublished(ctx context.Context, event events.TripPublished) error {
	args := m.Called(ctx, event)
	return args.Error(0)
}
