package endpoints

import (
	"context"
	"prueba_tecnica/api/entities"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateEvent(t *testing.T) {
	mockService := new(MockEventService)
	endpoints := NewEventEndpoints(mockService)
	ctx := context.Background()
	event := entities.Event{ID: "1", Name: "Test Event"}
	mockService.On("CreateEvent", ctx, event).Return(event, nil)

	result, err := endpoints.CreateEvent(ctx, event)

	assert.NoError(t, err)
	assert.Equal(t, event, result)
	mockService.AssertExpectations(t)
}

func TestGetEventByID(t *testing.T) {
	mockService := new(MockEventService)
	endpoints := NewEventEndpoints(mockService)
	ctx := context.Background()

	event := entities.Event{ID: "1", Name: "Test Event"}
	mockService.On("GetEventByID", ctx, "1").Return(event, nil)

	result, err := endpoints.GetEventByID(ctx, "1")

	assert.NoError(t, err)
	assert.Equal(t, event, result)
	mockService.AssertExpectations(t)
}

func TestDeleteEvent(t *testing.T) {
	mockService := new(MockEventService)
	endpoints := NewEventEndpoints(mockService)
	ctx := context.Background()

	mockService.On("DeleteEvent", ctx, "1").Return(nil)

	err := endpoints.DeleteEvent(ctx, "1")

	assert.NoError(t, err)
	mockService.AssertExpectations(t)
}

func TestUpdateEvent(t *testing.T) {
	mockService := new(MockEventService)
	endpoints := NewEventEndpoints(mockService)
	ctx := context.Background()
	event := entities.Event{ID: "1", Name: "Updated Event"}
	mockService.On("UpdateEvent", ctx, event).Return(event, nil)

	result, err := endpoints.UpdateEvent(ctx, event)

	assert.NoError(t, err)
	assert.Equal(t, event, result)
	mockService.AssertExpectations(t)
}

func TestClassifyEvent(t *testing.T) {
	mockService := new(MockEventService)
	endpoints := NewEventEndpoints(mockService)
	ctx := context.Background()
	event := entities.Event{ID: "1", Name: "Classified Event"}
	mockService.On("ClassifyEvent", ctx, "1").Return(event, nil)

	result, err := endpoints.ClassifyEvent(ctx, "1")

	assert.NoError(t, err)
	assert.Equal(t, event, result)
	mockService.AssertExpectations(t)
}

func TestManualClassifyEvent(t *testing.T) {
	mockService := new(MockEventService)
	endpoints := NewEventEndpoints(mockService)
	ctx := context.Background()
	event := entities.Event{ID: "1", Name: "Manually Classified Event"}
	mockService.On("ManualClassifyEvent", ctx, "1", "urgent").Return(event, nil)

	result, err := endpoints.ManualClassifyEvent(ctx, "1", "urgent")

	assert.NoError(t, err)
	assert.Equal(t, event, result)
	mockService.AssertExpectations(t)
}
