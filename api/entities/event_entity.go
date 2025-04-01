package entities

import "time"

type Event struct {
	ID          string    `json:"id,omitempty" bson:"_id,omitempty"`
	Name        string    `json:"name" bson:"name" validate:"required"`
	Type        string    `json:"type" bson:"type" validate:"required"`
	Description string    `json:"description" bson:"description" validate:"required"`
	Date        time.Time `json:"date" bson:"date"`
	Status      string    `json:"status" bson:"status" validate:"required"`     // "Pendiente" o "Revisado"
	Category    string    `json:"category,omitempty" bson:"category,omitempty"` // "Requiere gestión" o "Sin gestión"
	NeedsAction bool      `json:"needs_action,omitempty" bson:"needs_action,omitempty"`
}
