package transports

import (
	"prueba_tecnica/api/entities"

	"github.com/stretchr/testify/mock"
)

type MockService struct {
	mock.Mock
}

func (m *MockService) CreateEvent(event entities.Event) error {
	args := m.Called(event)
	return args.Error(0)
}

func (m *MockService) GetEvent(id string) (entities.Event, error) {
	args := m.Called(id)
	return args.Get(0).(entities.Event), args.Error(1)
}
