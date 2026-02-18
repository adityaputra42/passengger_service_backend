package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

// Booking holds the schema definition for the Booking entity.
type Booking struct {
	ent.Schema
}

// Fields of the Booking.
func (Booking) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(uuid.New),

		field.String("booking_code").
			Unique(),

		field.Float("total_amount"),

		field.String("status").
			Default("pending"),

		field.Time("expires_at").
			Optional().
			Nillable(),

		field.Time("created_at").
			Default(time.Now),
	}

}

// Edges of the Booking.
func (Booking) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("user", User.Type).
			Ref("bookings").
			Unique().
			Required(),

		edge.From("flight", Flight.Type).
			Ref("bookings").
			Unique().
			Required(),

		edge.To("booking_passengers", BookingPassenger.Type),
		edge.To("payments", Payment.Type),
	}
}
