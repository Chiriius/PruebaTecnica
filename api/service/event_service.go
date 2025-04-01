package service

import (
	"context"
	"errors"
	"prueba_tecnica/api/entities"
	"prueba_tecnica/api/repository"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
)

type EventService interface {
	CreateEvent(ctx context.Context, event entities.Event) (entities.Event, error)
	GetEventByID(ctx context.Context, id string) (entities.Event, error)
	GetAllEvents(ctx context.Context) ([]entities.Event, error)
	GetEventsByStatus(ctx context.Context, status string) ([]entities.Event, error)
	GetEventsByCategory(ctx context.Context, category string) ([]entities.Event, error)
	GetEventsNeedingAction(ctx context.Context) ([]entities.Event, error)
	UpdateEvent(ctx context.Context, event entities.Event) (entities.Event, error)
	DeleteEvent(ctx context.Context, id string) error
	ClassifyEvent(ctx context.Context, id string) (entities.Event, error)
	ManualClassifyEvent(ctx context.Context, id string, category string) (entities.Event, error)
}

type eventService struct {
	repo     repository.EventRepository
	logger   logrus.FieldLogger
	validate *validator.Validate
}

func NewEventService(repo repository.EventRepository, logger logrus.FieldLogger) EventService {
	return &eventService{
		repo:     repo,
		logger:   logger,
		validate: validator.New(),
	}
}

func (s *eventService) CreateEvent(ctx context.Context, event entities.Event) (entities.Event, error) {
	if err := s.validate.Struct(event); err != nil {
		s.logger.Errorln("Layer: event_service", "Method: CreateEvent", "Error:", err)
		return entities.Event{}, ErrValidation
	}

	if event.Status != "Pendiente por revisar" && event.Status != "Revisado" {
		s.logger.Errorln("Layer: event_service", "Method: CreateEvent", "Error:", ErrStatus)
		return entities.Event{}, ErrStatus
	}

	event.Date = time.Now()
	return s.repo.CreateEvent(ctx, event)
}

func (s *eventService) GetEventByID(ctx context.Context, id string) (entities.Event, error) {

	return s.repo.GetEventByID(ctx, id)
}

func (s *eventService) GetAllEvents(ctx context.Context) ([]entities.Event, error) {
	return s.repo.GetAllEvents(ctx)
}

func (s *eventService) GetEventsByStatus(ctx context.Context, status string) ([]entities.Event, error) {
	if status != "Pendiente por revisar" && status != "Revisado" {
		s.logger.Errorln("Layer: event_service", "Method: GetEventsByStatus", "Error:", ErrStatus)
		return nil, errors.New("status debe ser 'Pendiente por revisar' o 'Revisado'")
	}
	return s.repo.GetEventsByStatus(ctx, status)
}

func (s *eventService) GetEventsByCategory(ctx context.Context, category string) ([]entities.Event, error) {
	if category != "Requiere gestión" && category != "Sin gestión" {
		s.logger.Errorln("Layer: event_service", "Method: GetEventsByCategory", "Error:", ErrTypeCategory)
		return nil, ErrTypeCategory
	}
	return s.repo.GetEventsByCategory(ctx, category)
}

func (s *eventService) GetEventsNeedingAction(ctx context.Context) ([]entities.Event, error) {
	return s.repo.GetEventsNeedingAction(ctx)
}

func (s *eventService) UpdateEvent(ctx context.Context, event entities.Event) (entities.Event, error) {
	if err := s.validate.Struct(event); err != nil {
		s.logger.Errorln("Layer: user_services", "Method: UpdateUser", "Error:", err)
		return entities.Event{}, ErrValidation
	}
	if event.Status != "Pendiente por revisar" && event.Status != "Revisado" {
		s.logger.Errorln("Layer: event_service", "Method: CreateEvent", "Error:", ErrStatus)
		return entities.Event{}, ErrStatus
	}

	_, err := s.repo.GetEventByID(ctx, event.ID)
	if err != nil {
		s.logger.Errorln("Layer: event_service", "Method: CreateEvent", "Error:", err)
		return entities.Event{}, ErrEventNotfound
	}

	if event.Status == "Revisado" && event.Category == "" {
		if _, err := s.ClassifyEvent(ctx, event.ID); err != nil {
			s.logger.Errorln("Layer: event_service", "Method: CreateEvent", "Error:", err)
			return entities.Event{}, err
		}
	}

	return s.repo.UpdateEvent(ctx, event)
}

func (s *eventService) DeleteEvent(ctx context.Context, id string) error {
	if id == "" {
		s.logger.Errorln("Layer: event_service", "Method: DeleteEvent", "Error:", ErrNoID)
		return ErrNoID
	}
	return s.repo.DeleteEvent(ctx, id)
}

func (s *eventService) ClassifyEvent(ctx context.Context, id string) (entities.Event, error) {
	event, err := s.repo.GetEventByID(ctx, id)
	if err != nil {
		s.logger.Errorln("Layer: event_service", "Method: ClassifyEvent", "Error:", err)
		return entities.Event{}, err
	}

	if event.Status != "Revisado" {
		s.logger.Errorln("Layer: event_service", "Method: ClassifyEvent", "Error:", ErrEventRevi)
		return entities.Event{}, ErrEventRevi
	}

	switch event.Type {
	case "Incidente", "Problema", "Emergencia", "Error", "Critico":
		event.Category = "Requiere gestión"
		event.NeedsAction = true
	case "Reunión", "Informe", "Actualización", "Notificación", "Consulta":
		event.Category = "Sin gestión"
		event.NeedsAction = false
	default:
		event.Category = "Sin gestión"
		event.NeedsAction = false
	}

	return s.repo.UpdateEvent(ctx, event)
}

func (s *eventService) ManualClassifyEvent(ctx context.Context, id string, category string) (entities.Event, error) {
	if category != "Requiere gestión" && category != "Sin gestión" {
		s.logger.Errorln("Layer: event_service", "Method: ClassifyEvent", "Error:", ErrCategory)
		return entities.Event{}, ErrCategory
	}

	event, err := s.repo.GetEventByID(ctx, id)
	if err != nil {
		s.logger.Errorln("Layer: event_service", "Method: ClassifyEvent", "Error:", err)
		return entities.Event{}, err
	}

	if event.Status != "Revisado" {
		s.logger.Errorln("Layer: event_service", "Method: ClassifyEvent", "Error:", ErrEventRevi)
		return entities.Event{}, ErrEventRevi
	}

	event.Category = category
	event.NeedsAction = (category == "Requiere gestión")

	return s.repo.UpdateEvent(ctx, event)
}
