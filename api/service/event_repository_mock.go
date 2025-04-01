package service

import (
	"context"
	"prueba_tecnica/api/entities"

	"github.com/stretchr/testify/mock"
)

type mockEventRepository struct {
	mock.Mock
}

func (m *mockEventRepository) CreateEvent(ctx context.Context, event entities.Event) (entities.Event, error) {
	args := m.Called(ctx, event)
	return args.Get(0).(entities.Event), args.Error(1)
}

func (m *mockEventRepository) GetEventByID(ctx context.Context, id string) (entities.Event, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(entities.Event), args.Error(1)
}

func (m *mockEventRepository) GetAllEvents(ctx context.Context) ([]entities.Event, error) {
	args := m.Called(ctx)
	return args.Get(0).([]entities.Event), args.Error(1)
}

func (m *mockEventRepository) GetEventsByStatus(ctx context.Context, status string) ([]entities.Event, error) {
	args := m.Called(ctx, status)
	return args.Get(0).([]entities.Event), args.Error(1)
}

func (m *mockEventRepository) GetEventsByCategory(ctx context.Context, category string) ([]entities.Event, error) {
	args := m.Called(ctx, category)
	return args.Get(0).([]entities.Event), args.Error(1)
}

func (m *mockEventRepository) GetEventsNeedingAction(ctx context.Context) ([]entities.Event, error) {
	args := m.Called(ctx)
	return args.Get(0).([]entities.Event), args.Error(1)
}

func (m *mockEventRepository) UpdateEvent(ctx context.Context, event entities.Event) (entities.Event, error) {
	args := m.Called(ctx, event)
	return args.Get(0).(entities.Event), args.Error(1)
}

func (m *mockEventRepository) DeleteEvent(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
