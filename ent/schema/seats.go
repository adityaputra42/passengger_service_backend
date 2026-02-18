package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

// Seats holds the schema definition for the Seats entity.
type Seats struct {
	ent.Schema
}

// Fields of the Seats.
func (Seats) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(uuid.New),

		field.String("seat_number"),
		field.String("seat_class").
			Default("economy"),
	}
}

// Edges of the Seats.
func (Seats) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("aircraft", Aircraft.Type).
			Ref("seats").
			Unique().
			Required(),
		edge.To("flight_instances", FlightSeat.Type),
	}
}
