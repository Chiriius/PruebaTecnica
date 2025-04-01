package endpoints

import (
	"context"
	"prueba_tecnica/api/entities"
	"prueba_tecnica/api/service"
)

type EventEndpoints struct {
	CreateEvent            func(ctx context.Context, event entities.Event) (entities.Event, error)
	GetEventByID           func(ctx context.Context, id string) (entities.Event, error)
	GetAllEvents           func(ctx context.Context) ([]entities.Event, error)
	GetEventsByStatus      func(ctx context.Context, status string) ([]entities.Event, error)
	GetEventsByCategory    func(ctx context.Context, category string) ([]entities.Event, error)
	GetEventsNeedingAction func(ctx context.Context) ([]entities.Event, error)
	UpdateEvent            func(ctx context.Context, event entities.Event) (entities.Event, error)
	DeleteEvent            func(ctx context.Context, id string) error
	ClassifyEvent          func(ctx context.Context, id string) (entities.Event, error)
	ManualClassifyEvent    func(ctx context.Context, id string, category string) (entities.Event, error)
}

func NewEventEndpoints(s service.EventService) EventEndpoints {
	return EventEndpoints{
		CreateEvent:            s.CreateEvent,
		GetEventByID:           s.GetEventByID,
		GetAllEvents:           s.GetAllEvents,
		GetEventsByStatus:      s.GetEventsByStatus,
		GetEventsByCategory:    s.GetEventsByCategory,
		GetEventsNeedingAction: s.GetEventsNeedingAction,
		UpdateEvent:            s.UpdateEvent,
		DeleteEvent:            s.DeleteEvent,
		ClassifyEvent:          s.ClassifyEvent,
		ManualClassifyEvent:    s.ManualClassifyEvent,
	}
}
