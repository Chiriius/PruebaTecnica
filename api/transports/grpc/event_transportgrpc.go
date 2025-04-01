package transport

import (
	"context"
	"prueba_tecnica/api/endpoints"
	"prueba_tecnica/api/entities"
	pb "prueba_tecnica/api/pb/event"
	"time"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type EventHandler struct {
	pb.UnimplementedEventServiceServer
	endpoints endpoints.EventEndpoints
	logger    logrus.FieldLogger
}

func NewEventHandler(endpoints endpoints.EventEndpoints, logger logrus.FieldLogger) *EventHandler {
	return &EventHandler{
		endpoints: endpoints,
		logger:    logger,
	}
}

// Convertir Event de protobuf a entities.Event
func protoToEntity(protoEvent *pb.Event) entities.Event {
	var date time.Time
	if protoEvent.Date != nil {
		date = protoEvent.Date.AsTime()
	} else {
		date = time.Now()
	}

	return entities.Event{
		ID:          protoEvent.Id,
		Name:        protoEvent.Title,
		Description: protoEvent.Description,
		Type:        protoEvent.Type,
		Status:      protoEvent.Status,
		Category:    protoEvent.Category,
		Date:        date,
		NeedsAction: protoEvent.NeedsAction,
	}
}

// Convertir de entities.Event a protobuf Event
func entityToProto(event entities.Event) *pb.Event {
	return &pb.Event{
		Id:          event.ID,
		Title:       event.Name,
		Description: event.Description,
		Type:        event.Type,
		Status:      event.Status,
		Category:    event.Category,
		Date:        timestamppb.New(event.Date),
		NeedsAction: event.NeedsAction,
	}
}

// Implementaciones de los métodos del servicio gRPC
func (h *EventHandler) CreateEvent(ctx context.Context, req *pb.Event) (*pb.EventResponse, error) {
	h.logger.Infoln("Layer: grpc_handler", "Method: CreateEvent", "Request received")

	entityEvent := protoToEntity(req)
	event, err := h.endpoints.CreateEvent(ctx, entityEvent)
	if err != nil {
		h.logger.Errorln("Layer: grpc_handler", "Method: CreateEvent", "Error:", err)
		return nil, status.Errorf(codes.InvalidArgument, "failed to create event: %v", err)
	}

	return &pb.EventResponse{
		Id:      event.ID,
		Message: "Event created successfully",
	}, nil
}

func (h *EventHandler) GetEventByID(ctx context.Context, req *pb.EventID) (*pb.Event, error) {
	h.logger.Infoln("Layer: grpc_handler", "Method: GetEventByID", "Request received for ID:", req.Id)

	event, err := h.endpoints.GetEventByID(ctx, req.Id)
	if err != nil {
		h.logger.Errorln("Layer: grpc_handler", "Method: GetEventByID", "Error:", err)
		return nil, status.Errorf(codes.NotFound, "event not found: %v", err)
	}

	return entityToProto(event), nil
}

func (h *EventHandler) GetAllEvents(ctx context.Context, req *pb.Empty) (*pb.EventList, error) {
	h.logger.Infoln("Layer: grpc_handler", "Method: GetAllEvents", "Request received")

	events, err := h.endpoints.GetAllEvents(ctx)
	if err != nil {
		h.logger.Errorln("Layer: grpc_handler", "Method: GetAllEvents", "Error:", err)
		return nil, status.Errorf(codes.Internal, "failed to get events: %v", err)
	}

	protoEvents := make([]*pb.Event, len(events))
	for i, event := range events {
		protoEvents[i] = entityToProto(event)
	}

	return &pb.EventList{Events: protoEvents}, nil
}

func (h *EventHandler) GetEventsByStatus(ctx context.Context, req *pb.StatusRequest) (*pb.EventList, error) {
	h.logger.Infoln("Layer: grpc_handler", "Method: GetEventsByStatus", "Request received for status:", req.Status)

	events, err := h.endpoints.GetEventsByStatus(ctx, req.Status)
	if err != nil {
		h.logger.Errorln("Layer: grpc_handler", "Method: GetEventsByStatus", "Error:", err)
		return nil, status.Errorf(codes.InvalidArgument, "failed to get events by status: %v", err)
	}

	protoEvents := make([]*pb.Event, len(events))
	for i, event := range events {
		protoEvents[i] = entityToProto(event)
	}

	return &pb.EventList{Events: protoEvents}, nil
}

func (h *EventHandler) GetEventsByCategory(ctx context.Context, req *pb.CategoryRequest) (*pb.EventList, error) {
	h.logger.Infoln("Layer: grpc_handler", "Method: GetEventsByCategory", "Request received for category:", req.Category)

	events, err := h.endpoints.GetEventsByCategory(ctx, req.Category)
	if err != nil {
		h.logger.Errorln("Layer: grpc_handler", "Method: GetEventsByCategory", "Error:", err)
		return nil, status.Errorf(codes.InvalidArgument, "failed to get events by category: %v", err)
	}

	protoEvents := make([]*pb.Event, len(events))
	for i, event := range events {
		protoEvents[i] = entityToProto(event)
	}

	return &pb.EventList{Events: protoEvents}, nil
}

func (h *EventHandler) GetEventsNeedingAction(ctx context.Context, req *pb.Empty) (*pb.EventList, error) {
	h.logger.Infoln("Layer: grpc_handler", "Method: GetEventsNeedingAction", "Request received")

	events, err := h.endpoints.GetEventsNeedingAction(ctx)
	if err != nil {
		h.logger.Errorln("Layer: grpc_handler", "Method: GetEventsNeedingAction", "Error:", err)
		return nil, status.Errorf(codes.Internal, "failed to get events needing action: %v", err)
	}

	protoEvents := make([]*pb.Event, len(events))
	for i, event := range events {
		protoEvents[i] = entityToProto(event)
	}

	return &pb.EventList{Events: protoEvents}, nil
}

func (h *EventHandler) UpdateEvent(ctx context.Context, req *pb.Event) (*pb.Event, error) {
	h.logger.Infoln("Layer: grpc_handler", "Method: UpdateEvent", "Request received for ID:", req.Id)

	entityEvent := protoToEntity(req)
	event, err := h.endpoints.UpdateEvent(ctx, entityEvent)
	if err != nil {
		h.logger.Errorln("Layer: grpc_handler", "Method: UpdateEvent", "Error:", err)
		return nil, status.Errorf(codes.InvalidArgument, "failed to update event: %v", err)
	}

	return entityToProto(event), nil
}

func (h *EventHandler) DeleteEvent(ctx context.Context, req *pb.EventID) (*pb.DeleteResponse, error) {
	h.logger.Infoln("Layer: grpc_handler", "Method: DeleteEvent", "Request received for ID:", req.Id)

	err := h.endpoints.DeleteEvent(ctx, req.Id)
	if err != nil {
		h.logger.Errorln("Layer: grpc_handler", "Method: DeleteEvent", "Error:", err)
		return nil, status.Errorf(codes.InvalidArgument, "failed to delete event: %v", err)
	}

	return &pb.DeleteResponse{
		Success: true,
		Message: "Event deleted successfully",
	}, nil
}

func (h *EventHandler) ClassifyEvent(ctx context.Context, req *pb.EventID) (*pb.Event, error) {
	h.logger.Infoln("Layer: grpc_handler", "Method: ClassifyEvent", "Request received for ID:", req.Id)

	event, err := h.endpoints.ClassifyEvent(ctx, req.Id)
	if err != nil {
		h.logger.Errorln("Layer: grpc_handler", "Method: ClassifyEvent", "Error:", err)
		return nil, status.Errorf(codes.InvalidArgument, "failed to classify event: %v", err)
	}

	return entityToProto(event), nil
}

func (h *EventHandler) ManualClassifyEvent(ctx context.Context, req *pb.ManualClassifyRequest) (*pb.Event, error) {
	h.logger.Infoln("Layer: grpc_handler", "Method: ManualClassifyEvent", "Request received for ID:", req.Id)

	event, err := h.endpoints.ManualClassifyEvent(ctx, req.Id, req.Category)
	if err != nil {
		h.logger.Errorln("Layer: grpc_handler", "Method: ManualClassifyEvent", "Error:", err)
		return nil, status.Errorf(codes.InvalidArgument, "failed to manually classify event: %v", err)
	}

	return entityToProto(event), nil
}
