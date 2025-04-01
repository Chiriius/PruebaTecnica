package repository

import (
	"context"
	"prueba_tecnica/api/entities"
	"time"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type EventRepository interface {
	CreateEvent(ctx context.Context, event entities.Event) (entities.Event, error)
	GetEventByID(ctx context.Context, id string) (entities.Event, error)
	GetAllEvents(ctx context.Context) ([]entities.Event, error)
	GetEventsByStatus(ctx context.Context, status string) ([]entities.Event, error)
	GetEventsByCategory(ctx context.Context, category string) ([]entities.Event, error)
	GetEventsNeedingAction(ctx context.Context) ([]entities.Event, error)
	UpdateEvent(ctx context.Context, event entities.Event) (entities.Event, error)
	DeleteEvent(ctx context.Context, id string) error
}

type MongoEventRepository struct {
	db     *mongo.Client
	logger logrus.FieldLogger
}

func NewMongoEventRepository(db *mongo.Client, logger logrus.FieldLogger) *MongoEventRepository {
	return &MongoEventRepository{
		db:     db,
		logger: logger,
	}

}

func (r *MongoEventRepository) CreateEvent(ctx context.Context, event entities.Event) (entities.Event, error) {
	coll := r.db.Database("events_db").Collection("events")
	event.Date = time.Now()
	result, err := coll.InsertOne(ctx, event)

	if err != nil {
		r.logger.Errorln("Layer:event_repository ", "Method:CreateEvent ", "Error:", err)
		return event, err
	}

	event.ID = result.InsertedID.(primitive.ObjectID).Hex()
	r.logger.Infoln("Layer:event_repository", "Method:CreateEvent", "event:", event)
	return event, err
}

func (r *MongoEventRepository) GetEventByID(ctx context.Context, id string) (entities.Event, error) {
	var event entities.Event
	idd, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		r.logger.Errorln("Layer:event_repository ", "Method:GetEventByID ", "Error:", err)
		return event, ErrEventNotfound
	}

	filter := bson.D{{"_id", idd}}
	opts := options.FindOne()
	coll := r.db.Database("events_db").Collection("events")

	err = coll.FindOne(ctx, filter, opts).Decode(&event)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			r.logger.Errorln("Layer:event_repository ", "Method:GetEventByID ", "Error:", ErrEventNotfound)
			return event, ErrEventNotfound
		}
		r.logger.Errorln("Layer:event_repository ", "Method:GetEventByID ", "Error:", err)
		return event, err
	}
	r.logger.Infoln("Layer:event_repository", "Method:GetEventByID", "event:", event)
	return event, nil

}

func (r *MongoEventRepository) GetAllEvents(ctx context.Context) ([]entities.Event, error) {
	return r.findEvents(ctx, bson.M{})
}

func (r *MongoEventRepository) GetEventsByStatus(ctx context.Context, status string) ([]entities.Event, error) {
	r.logger.Infoln("Layer:event_repository", "Method:GetEventsByStatus", "status:", status)
	return r.findEvents(ctx, bson.M{"status": status})
}

func (r *MongoEventRepository) GetEventsByCategory(ctx context.Context, category string) ([]entities.Event, error) {
	return r.findEvents(ctx, bson.M{
		"category": category,
		"status":   "Revisado",
	})
}

func (r *MongoEventRepository) GetEventsNeedingAction(ctx context.Context) ([]entities.Event, error) {
	return r.findEvents(ctx, bson.M{
		"needs_action": true,
		"status":       "Revisado",
	})
}

func (r *MongoEventRepository) UpdateEvent(ctx context.Context, event entities.Event) (entities.Event, error) {
	ide := string(event.ID)
	idd, err := primitive.ObjectIDFromHex(ide)
	if err != nil {
		r.logger.Errorln("Layer:event_repository", "Method:UpdateEvent ", "Error:", err)
		return entities.Event{}, err
	}

	coll := r.db.Database("events_db").Collection("events")

	filter := bson.D{{"_id", idd}}
	update := bson.M{
		"$set": bson.M{
			"name":        event.Name,
			"type":        event.Type,
			"description": event.Description,
			"date":        event.Date,
			"status":      event.Status,
		},
	}

	_, err = coll.UpdateOne(ctx, filter, update)
	if err != nil {
		r.logger.Errorln("Layer:event_repository ", "Method:UpdateEvent ", "Error:", err)
		return entities.Event{}, err
	}

	r.logger.Infoln("Layer:event_repository ", "Method:UpdateEvent ", "Evento Actualizado:", event)
	return event, err
}

func (r *MongoEventRepository) DeleteEvent(ctx context.Context, id string) error {
	coll := r.db.Database("events_db").Collection("events")
	idd, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		r.logger.Errorln("Layer:event_repository ", "Method: DeleteEvent ", "Error:", err)
		return err
	}

	filter := bson.D{{"_id", idd}}
	res, err := coll.DeleteOne(ctx, filter)
	if err != nil {
		r.logger.Errorln("Layer:event_repository ", "Method: DeleteEvent ", "Error:", err)

		return ErrEventNotfound
	}

	if res.DeletedCount == 0 {
		r.logger.Errorln("Layer:event_repository ", "Method: DeleteEvent ", "Error: No tasks were deleted")
		return ErrNotasks
	}
	r.logger.Infoln("Layer:event_repository ", "Method: DeleteEvent ", "Event:", idd)
	return err
}

func (r *MongoEventRepository) findEvents(ctx context.Context, filter bson.M) ([]entities.Event, error) {
	coll := r.db.Database("events_db").Collection("events")
	cursor, err := coll.Find(ctx, filter, options.Find().SetSort(bson.M{"date": -1}))
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var events []entities.Event
	for cursor.Next(ctx) {
		var event entities.Event
		if err := cursor.Decode(&event); err != nil {
			r.logger.Errorln("Layer:event_repository ", "Method:findEvents ", "Error:", err)
			return nil, err
		}
		events = append(events, event)
	}
	if err := cursor.Err(); err != nil {
		r.logger.Errorln("Layer:event_repository ", "Method:findEvents ", "Error:", err)
		return nil, err
	}
	r.logger.Infoln("Layer:event_repository", "Method:findEvents", "eventos econtrados correctamente")
	return events, nil
}
