package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

// Flight holds the schema definition for the Flight entity.
type Flight struct {
	ent.Schema
}

// Fields of the Flight.
func (Flight) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(uuid.New),

		field.String("flight_number"),

		field.Time("departure_time"),
		field.Time("arrival_time"),

		field.Float("base_price"),

		field.String("status").
			Default("scheduled"),

		field.Time("created_at").
			Default(time.Now),
	}
}

// Edges of the Flight.
func (Flight) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("aircraft", Aircraft.Type).
			Ref("flights").
			Unique().
			Required(),

		edge.From("departure_airport", Airport.Type).
			Ref("departures").
			Unique().
			Required(),

		edge.From("arrival_airport", Airport.Type).
			Ref("arrivals").
			Unique().
			Required(),

		edge.To("flight_seats", FlightSeat.Type),
		edge.To("bookings", Booking.Type),
	}
}
