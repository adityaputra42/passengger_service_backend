package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

// FlightSeat holds the schema definition for the FlightSeat entity.
type FlightSeat struct {
	ent.Schema
}

// Fields of the FlightSeat.
func (FlightSeat) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(uuid.New),

		field.String("status").
			Default("available"),

		field.Float("price"),
	}
}

// Edges of the FlightSeat.
func (FlightSeat) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("flight", Flight.Type).
			Ref("flight_seats").
			Unique().
			Required(),

		edge.From("seat", Seats.Type).
			Ref("flight_instances").
			Unique().
			Required(),
		edge.To("assigned", BookingPassenger.Type),
	}
}
