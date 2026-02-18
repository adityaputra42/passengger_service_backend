package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

// BookingPassenger holds the schema definition for the BookingPassenger entity.
type BookingPassenger struct {
	ent.Schema
}

// Fields of the BookingPassenger.
func (BookingPassenger) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(uuid.New),
	}
}

// Edges of the BookingPassenger.
func (BookingPassenger) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("booking", Booking.Type).
			Ref("booking_passengers").
			Unique().
			Required(),

		edge.From("passenger", Passenger.Type).
			Ref("booking_entries").
			Unique().
			Required(),

		edge.From("flight_seat", FlightSeat.Type).
			Ref("assigned").
			Unique().
			Required(),
	}
}
