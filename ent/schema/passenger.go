package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

// Passenger holds the schema definition for the Passenger entity.
type Passenger struct {
	ent.Schema
}

// Fields of the Passenger.
func (Passenger) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(uuid.New),

		field.String("first_name"),
		field.String("last_name"),
		field.Time("date_of_birth"),

		field.String("passport_number").
			Optional().
			Nillable(),

		field.Time("created_at").
			Default(time.Now),
	}
}

// Edges of the Passenger.
func (Passenger) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("booking_entries", BookingPassenger.Type),
	}
}
