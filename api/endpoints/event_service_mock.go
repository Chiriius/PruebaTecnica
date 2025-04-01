package endpoints

import (
	"context"
	"prueba_tecnica/api/entities"

	"github.com/stretchr/testify/mock"
)

type MockEventService struct {
	mock.Mock
}

func (m *MockEventService) CreateEvent(ctx context.Context, event entities.Event) (entities.Event, error) {
	args := m.Called(ctx, event)
	return args.Get(0).(entities.Event), args.Error(1)
}

func (m *MockEventService) GetEventByID(ctx context.Context, id string) (entities.Event, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(entities.Event), args.Error(1)
}

func (m *MockEventService) GetAllEvents(ctx context.Context) ([]entities.Event, error) {
	args := m.Called(ctx)
	return args.Get(0).([]entities.Event), args.Error(1)
}

func (m *MockEventService) GetEventsByStatus(ctx context.Context, status string) ([]entities.Event, error) {
	args := m.Called(ctx, status)
	return args.Get(0).([]entities.Event), args.Error(1)
}

func (m *MockEventService) GetEventsByCategory(ctx context.Context, category string) ([]entities.Event, error) {
	args := m.Called(ctx, category)
	return args.Get(0).([]entities.Event), args.Error(1)
}

func (m *MockEventService) GetEventsNeedingAction(ctx context.Context) ([]entities.Event, error) {
	args := m.Called(ctx)
	return args.Get(0).([]entities.Event), args.Error(1)
}

func (m *MockEventService) UpdateEvent(ctx context.Context, event entities.Event) (entities.Event, error) {
	args := m.Called(ctx, event)
	return args.Get(0).(entities.Event), args.Error(1)
}

func (m *MockEventService) DeleteEvent(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockEventService) ClassifyEvent(ctx context.Context, id string) (entities.Event, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(entities.Event), args.Error(1)
}

func (m *MockEventService) ManualClassifyEvent(ctx context.Context, id string, category string) (entities.Event, error) {
	args := m.Called(ctx, id, category)
	return args.Get(0).(entities.Event), args.Error(1)
}
