package service

import (
	"context"
	"errors"
	"prueba_tecnica/api/entities"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateEvent(t *testing.T) {
	testCases := []struct {
		name          string
		event         entities.Event
		mockResponse  entities.Event
		mockError     error
		expectedError error
	}{
		{
			name: "Success - Create event with valid data",
			event: entities.Event{
				Name:        "Test Event",
				Description: "Test Description",
				Type:        "Incidente",
				Status:      "Pendiente por revisar",
			},
			mockResponse: entities.Event{
				ID:          "1",
				Name:        "Test Event",
				Description: "Test Description",
				Type:        "Incidente",
				Status:      "Pendiente por revisar",
				Date:        time.Now(),
			},
			mockError:     nil,
			expectedError: nil,
		},
		{
			name: "Failure - Invalid status",
			event: entities.Event{
				Name:        "Test Event",
				Description: "Test Description",
				Type:        "Incidente",
				Status:      "Invalid Status",
			},
			mockResponse:  entities.Event{},
			mockError:     nil,
			expectedError: ErrStatus,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			mockRepo := new(mockEventRepository)
			logger := logrus.New()
			service := NewEventService(mockRepo, logger)

			if tc.expectedError == nil || tc.expectedError != ErrStatus {
				mockRepo.On("CreateEvent", mock.Anything, mock.AnythingOfType("entities.Event")).Return(tc.mockResponse, tc.mockError)
			}

			// Execute
			result, err := service.CreateEvent(context.Background(), tc.event)

			// Assert
			if tc.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tc.expectedError, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.mockResponse.ID, result.ID)
				mockRepo.AssertExpectations(t)
			}
		})
	}
}

func TestGetEventByID(t *testing.T) {
	testCases := []struct {
		name          string
		eventID       string
		mockResponse  entities.Event
		mockError     error
		expectedError error
	}{
		{
			name:    "Success - Get event by ID",
			eventID: "1",
			mockResponse: entities.Event{
				ID:          "1",
				Name:        "Test Event",
				Description: "Test Description",
				Type:        "Incidente",
				Status:      "Pendiente por revisar",
			},
			mockError:     nil,
			expectedError: nil,
		},
		{
			name:          "Failure - Event not found",
			eventID:       "999",
			mockResponse:  entities.Event{},
			mockError:     ErrEventNotfound,
			expectedError: ErrEventNotfound,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Setup
			mockRepo := new(mockEventRepository)
			logger := logrus.New()
			service := NewEventService(mockRepo, logger)

			mockRepo.On("GetEventByID", mock.Anything, tc.eventID).Return(tc.mockResponse, tc.mockError)

			// Execute
			result, err := service.GetEventByID(context.Background(), tc.eventID)

			// Assert
			if tc.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tc.expectedError, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.mockResponse.ID, result.ID)
			}
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestGetAllEvents(t *testing.T) {
	mockEvents := []entities.Event{
		{
			ID:          "1",
			Name:        "Event 1",
			Description: "Description 1",
			Type:        "Incidente",
			Status:      "Pendiente por revisar",
		},
		{
			ID:          "2",
			Name:        "Event 2",
			Description: "Description 2",
			Type:        "Problema",
			Status:      "Revisado",
		},
	}

	testCases := []struct {
		name          string
		mockResponse  []entities.Event
		mockError     error
		expectedError error
		expectedCount int
	}{
		{
			name:          "Success - Get all events",
			mockResponse:  mockEvents,
			mockError:     nil,
			expectedError: nil,
			expectedCount: 2,
		},
		{
			name:          "Failure - Repository error",
			mockResponse:  []entities.Event{},
			mockError:     errors.New("database error"),
			expectedError: errors.New("database error"),
			expectedCount: 0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			mockRepo := new(mockEventRepository)
			logger := logrus.New()
			service := NewEventService(mockRepo, logger)

			mockRepo.On("GetAllEvents", mock.Anything).Return(tc.mockResponse, tc.mockError)

			// Execute
			result, err := service.GetAllEvents(context.Background())

			// Assert
			if tc.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tc.expectedError.Error(), err.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedCount, len(result))
			}
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestGetEventsByStatus(t *testing.T) {
	mockEvents := []entities.Event{
		{
			ID:          "1",
			Name:        "Event 1",
			Description: "Description 1",
			Type:        "Incidente",
			Status:      "Pendiente por revisar",
		},
		{
			ID:          "2",
			Name:        "Event 2",
			Description: "Description 2",
			Type:        "Problema",
			Status:      "Pendiente por revisar",
		},
	}

	testCases := []struct {
		name          string
		status        string
		mockResponse  []entities.Event
		mockError     error
		expectedError error
		expectedCount int
	}{
		{
			name:          "Success - Get events by valid status",
			status:        "Pendiente por revisar",
			mockResponse:  mockEvents,
			mockError:     nil,
			expectedError: nil,
			expectedCount: 2,
		},
		{
			name:          "Failure - Invalid status",
			status:        "Invalid Status",
			mockResponse:  []entities.Event{},
			mockError:     nil,
			expectedError: errors.New("status debe ser 'Pendiente por revisar' o 'Revisado'"),
			expectedCount: 0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			mockRepo := new(mockEventRepository)
			logger := logrus.New()
			service := NewEventService(mockRepo, logger)

			if tc.expectedError == nil {
				mockRepo.On("GetEventsByStatus", mock.Anything, tc.status).Return(tc.mockResponse, tc.mockError)
			}

			// Execute
			result, err := service.GetEventsByStatus(context.Background(), tc.status)

			// Assert
			if tc.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tc.expectedError.Error(), err.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedCount, len(result))
			}
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestGetEventsByCategory(t *testing.T) {
	mockEvents := []entities.Event{
		{
			ID:          "1",
			Name:        "Event 1",
			Description: "Description 1",
			Type:        "Incidente",
			Status:      "Revisado",
			Category:    "Requiere gestión",
		},
	}

	testCases := []struct {
		name          string
		category      string
		mockResponse  []entities.Event
		mockError     error
		expectedError error
		expectedCount int
	}{
		{
			name:          "Success - Get events by valid category",
			category:      "Requiere gestión",
			mockResponse:  mockEvents,
			mockError:     nil,
			expectedError: nil,
			expectedCount: 1,
		},
		{
			name:          "Failure - Invalid category",
			category:      "Invalid Category",
			mockResponse:  []entities.Event{},
			mockError:     nil,
			expectedError: ErrTypeCategory,
			expectedCount: 0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			mockRepo := new(mockEventRepository)
			logger := logrus.New()
			service := NewEventService(mockRepo, logger)

			if tc.expectedError == nil {
				mockRepo.On("GetEventsByCategory", mock.Anything, tc.category).Return(tc.mockResponse, tc.mockError)
			}

			// Execute
			result, err := service.GetEventsByCategory(context.Background(), tc.category)

			// Assert
			if tc.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tc.expectedError, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedCount, len(result))
			}
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestGetEventsNeedingAction(t *testing.T) {
	mockEvents := []entities.Event{
		{
			ID:          "1",
			Name:        "Event 1",
			Description: "Description 1",
			Type:        "Incidente",
			Status:      "Revisado",
			Category:    "Requiere gestión",
			NeedsAction: true,
		},
	}

	testCases := []struct {
		name          string
		mockResponse  []entities.Event
		mockError     error
		expectedError error
		expectedCount int
	}{
		{
			name:          "Success - Get events needing action",
			mockResponse:  mockEvents,
			mockError:     nil,
			expectedError: nil,
			expectedCount: 1,
		},
		{
			name:          "Failure - Repository error",
			mockResponse:  []entities.Event{},
			mockError:     errors.New("database error"),
			expectedError: errors.New("database error"),
			expectedCount: 0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Setup
			mockRepo := new(mockEventRepository)
			logger := logrus.New()
			service := NewEventService(mockRepo, logger)

			mockRepo.On("GetEventsNeedingAction", mock.Anything).Return(tc.mockResponse, tc.mockError)

			// Execute
			result, err := service.GetEventsNeedingAction(context.Background())

			// Assert
			if tc.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tc.expectedError.Error(), err.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedCount, len(result))
			}
			mockRepo.AssertExpectations(t)
		})
	}
}

/*
	func TestUpdateEvent(t *testing.T) {
		existingEvent := entities.Event{
			ID:          "1",
			Name:        "Old Title",
			Description: "Old Description",
			Type:        "Incidente",
			Status:      "Pendiente por revisar",
		}

		updatedEvent := entities.Event{
			ID:          "1",
			Name:        "Updated Title",
			Description: "Updated Description",
			Type:        "Incidente",
			Status:      "Revisado",
		}

		classifiedEvent := entities.Event{
			ID:          "1",
			Name:        "Updated Title",
			Description: "Updated Description",
			Type:        "Incidente",
			Status:      "Revisado",
			Category:    "Requiere gestión",
			NeedsAction: true,
		}

		testCases := []struct {
			name           string
			event          entities.Event
			getEventReturn entities.Event
			getEventError  error
			updateReturn   entities.Event
			updateError    error
			expectedError  error
		}{
			{
				name:           "Success - Update event",
				event:          updatedEvent,
				getEventReturn: existingEvent,
				getEventError:  nil,
				updateReturn:   updatedEvent,
				updateError:    nil,
				expectedError:  nil,
			},
			{
				name: "Failure - Invalid status",
				event: entities.Event{
					ID:          "1",
					Name:        "Updated Title",
					Description: "Updated Description",
					Type:        "Incidente",
					Status:      "Invalid Status",
				},
				getEventReturn: entities.Event{},
				getEventError:  nil,
				updateReturn:   entities.Event{},
				updateError:    nil,
				expectedError:  ErrStatus,
			},
			{
				name:           "Failure - Event not found",
				event:          updatedEvent,
				getEventReturn: entities.Event{},
				getEventError:  ErrEventNotfound,
				updateReturn:   entities.Event{},
				updateError:    nil,
				expectedError:  ErrEventNotfound,
			},
			{
				name: "Success - Auto classify event when status is Revisado",
				event: entities.Event{
					ID:          "1",
					Name:        "Updated Title",
					Description: "Updated Description",
					Type:        "Incidente",
					Status:      "Revisado",
				},
				getEventReturn: existingEvent,
				getEventError:  nil,
				updateReturn:   classifiedEvent,
				updateError:    nil,
				expectedError:  nil,
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				// Setup
				mockRepo := new(mockEventRepository)
				logger := logrus.New()
				service := NewEventService(mockRepo, logger)

				if tc.expectedError != ErrStatus {
					mockRepo.On("GetEventByID", mock.Anything, tc.event.ID).Return(tc.getEventReturn, tc.getEventError)

					// For testing auto-classification
					if tc.event.Status == "Revisado" && tc.event.Category == "" && tc.getEventError == nil {
						classEvent := tc.event
						classEvent.Category = "Requiere gestión"
						classEvent.NeedsAction = true
						mockRepo.On("UpdateEvent", mock.Anything, mock.MatchedBy(func(e entities.Event) bool {
							return e.ID == tc.event.ID && e.Category == "Requiere gestión"
						})).Return(classifiedEvent, nil)
					}

					if tc.getEventError == nil {
						mockRepo.On("UpdateEvent", mock.Anything, tc.event).Return(tc.updateReturn, tc.updateError)
					}
				}

				// Execute
				result, err := service.UpdateEvent(context.Background(), tc.event)

				// Assert
				if tc.expectedError != nil {
					assert.Error(t, err)
					if tc.expectedError.Error() != err.Error() {
						t.Errorf("Expected error %v, got %v", tc.expectedError, err)
					}
				} else {
					assert.NoError(t, err)
					assert.Equal(t, tc.event.ID, result.ID)
				}

				mockRepo.AssertExpectations(t)
			})
		}
	}
*/
func TestDeleteEvent(t *testing.T) {
	testCases := []struct {
		name          string
		eventID       string
		mockError     error
		expectedError error
	}{
		{
			name:          "Success - Delete event",
			eventID:       "1",
			mockError:     nil,
			expectedError: nil,
		},
		{
			name:          "Failure - Empty ID",
			eventID:       "",
			mockError:     nil,
			expectedError: ErrNoID,
		},
		{
			name:          "Failure - Repository error",
			eventID:       "1",
			mockError:     errors.New("database error"),
			expectedError: errors.New("database error"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Setup
			mockRepo := new(mockEventRepository)
			logger := logrus.New()
			service := NewEventService(mockRepo, logger)

			if tc.eventID != "" {
				mockRepo.On("DeleteEvent", mock.Anything, tc.eventID).Return(tc.mockError)
			}

			// Execute
			err := service.DeleteEvent(context.Background(), tc.eventID)

			// Assert
			if tc.expectedError != nil {
				assert.Error(t, err)
				if tc.expectedError.Error() != err.Error() {
					t.Errorf("Expected error %v, got %v", tc.expectedError, err)
				}
			} else {
				assert.NoError(t, err)
			}
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestClassifyEvent(t *testing.T) {
	testCases := []struct {
		name             string
		eventID          string
		getEventReturn   entities.Event
		getEventError    error
		updateReturn     entities.Event
		updateError      error
		expectedCategory string
		expectedAction   bool
		expectedError    error
	}{
		{
			name:    "Success - Classify event requiring action",
			eventID: "1",
			getEventReturn: entities.Event{
				ID:     "1",
				Name:   "Test Event",
				Type:   "Incidente",
				Status: "Revisado",
			},
			getEventError: nil,
			updateReturn: entities.Event{
				ID:          "1",
				Name:        "Test Event",
				Type:        "Incidente",
				Status:      "Revisado",
				Category:    "Requiere gestión",
				NeedsAction: true,
			},
			updateError:      nil,
			expectedCategory: "Requiere gestión",
			expectedAction:   true,
			expectedError:    nil,
		},
		{
			name:    "Success - Classify event not requiring action",
			eventID: "2",
			getEventReturn: entities.Event{
				ID:     "2",
				Name:   "Test Event",
				Type:   "Reunión",
				Status: "Revisado",
			},
			getEventError: nil,
			updateReturn: entities.Event{
				ID:          "2",
				Name:        "Test Event",
				Type:        "Reunión",
				Status:      "Revisado",
				Category:    "Sin gestión",
				NeedsAction: false,
			},
			updateError:      nil,
			expectedCategory: "Sin gestión",
			expectedAction:   false,
			expectedError:    nil,
		},
		{
			name:           "Failure - Event not found",
			eventID:        "999",
			getEventReturn: entities.Event{},
			getEventError:  ErrEventNotfound,
			updateReturn:   entities.Event{},
			updateError:    nil,
			expectedError:  ErrEventNotfound,
		},
		{
			name:    "Failure - Event not reviewed",
			eventID: "3",
			getEventReturn: entities.Event{
				ID:     "3",
				Name:   "Test Event",
				Type:   "Incidente",
				Status: "Pendiente por revisar",
			},
			getEventError: nil,
			updateReturn:  entities.Event{},
			updateError:   nil,
			expectedError: ErrEventRevi,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Setup
			mockRepo := new(mockEventRepository)
			logger := logrus.New()
			service := NewEventService(mockRepo, logger)

			mockRepo.On("GetEventByID", mock.Anything, tc.eventID).Return(tc.getEventReturn, tc.getEventError)

			if tc.getEventError == nil && tc.getEventReturn.Status == "Revisado" {
				expectedEvent := tc.getEventReturn
				expectedEvent.Category = tc.expectedCategory
				expectedEvent.NeedsAction = tc.expectedAction

				mockRepo.On("UpdateEvent", mock.Anything, mock.MatchedBy(func(e entities.Event) bool {
					return e.ID == tc.eventID && e.Category == tc.expectedCategory && e.NeedsAction == tc.expectedAction
				})).Return(tc.updateReturn, tc.updateError)
			}

			// Execute
			result, err := service.ClassifyEvent(context.Background(), tc.eventID)

			// Assert
			if tc.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tc.expectedError, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedCategory, result.Category)
				assert.Equal(t, tc.expectedAction, result.NeedsAction)
			}
			mockRepo.AssertExpectations(t)
		})
	}
}
func TestManualClassifyEvent(t *testing.T) {
	testCases := []struct {
		name           string
		eventID        string
		category       string
		getEventReturn entities.Event
		getEventError  error
		updateReturn   entities.Event
		updateError    error
		expectedAction bool
		expectedError  error
	}{
		{
			name:     "Success - Manual classify requiring action",
			eventID:  "1",
			category: "Requiere gestión",
			getEventReturn: entities.Event{
				ID:     "1",
				Name:   "Test Event",
				Type:   "Incidente",
				Status: "Revisado",
			},
			getEventError: nil,
			updateReturn: entities.Event{
				ID:          "1",
				Name:        "Test Event",
				Type:        "Incidente",
				Status:      "Revisado",
				Category:    "Requiere gestión",
				NeedsAction: true,
			},
			updateError:    nil,
			expectedAction: true,
			expectedError:  nil,
		},
		{
			name:     "Success - Manual classify not requiring action",
			eventID:  "2",
			category: "Sin gestión",
			getEventReturn: entities.Event{
				ID:     "2",
				Name:   "Test Event",
				Type:   "Incidente",
				Status: "Revisado",
			},
			getEventError: nil,
			updateReturn: entities.Event{
				ID:          "2",
				Name:        "Test Event",
				Type:        "Incidente",
				Status:      "Revisado",
				Category:    "Sin gestión",
				NeedsAction: false,
			},
			updateError:    nil,
			expectedAction: false,
			expectedError:  nil,
		},
		{
			name:           "Failure - Invalid category",
			eventID:        "1",
			category:       "Invalid Category",
			getEventReturn: entities.Event{},
			getEventError:  nil,
			updateReturn:   entities.Event{},
			updateError:    nil,
			expectedError:  ErrCategory,
		},
		{
			name:           "Failure - Event not found",
			eventID:        "999",
			category:       "Requiere gestión",
			getEventReturn: entities.Event{},
			getEventError:  ErrEventNotfound,
			updateReturn:   entities.Event{},
			updateError:    nil,
			expectedError:  ErrEventNotfound,
		},
		{
			name:     "Failure - Event not reviewed",
			eventID:  "3",
			category: "Requiere gestión",
			getEventReturn: entities.Event{
				ID:     "3",
				Name:   "Test Event",
				Type:   "Incidente",
				Status: "Pendiente por revisar",
			},
			getEventError: nil,
			updateReturn:  entities.Event{},
			updateError:   nil,
			expectedError: ErrEventRevi,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Setup
			mockRepo := new(mockEventRepository)
			logger := logrus.New()
			service := NewEventService(mockRepo, logger)

			if tc.category == "Requiere gestión" || tc.category == "Sin gestión" {
				mockRepo.On("GetEventByID", mock.Anything, tc.eventID).Return(tc.getEventReturn, tc.getEventError)

				if tc.getEventError == nil && tc.getEventReturn.Status == "Revisado" {
					expectedEvent := tc.getEventReturn
					expectedEvent.Category = tc.category
					expectedEvent.NeedsAction = tc.expectedAction

					mockRepo.On("UpdateEvent", mock.Anything, mock.MatchedBy(func(e entities.Event) bool {
						return e.ID == tc.eventID && e.Category == tc.category && e.NeedsAction == tc.expectedAction
					})).Return(tc.updateReturn, tc.updateError)
				}
			}

			// Execute
			result, err := service.ManualClassifyEvent(context.Background(), tc.eventID, tc.category)

			// Assert
			if tc.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tc.expectedError, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.category, result.Category)
				assert.Equal(t, tc.expectedAction, result.NeedsAction)
			}
			mockRepo.AssertExpectations(t)
		})
	}
}
