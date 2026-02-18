package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

// Payment holds the schema definition for the Payment entity.
type Payment struct {
	ent.Schema
}

// Fields of the Payment.
func (Payment) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(uuid.New),

		field.Float("amount"),

		field.String("provider"),

		field.String("status").
			Default("pending"),

		field.Time("paid_at").
			Optional().
			Nillable(),

		field.Time("created_at").
			Default(time.Now),
	}
}

// Edges of the Payment.
func (Payment) Edges() []ent.Edge {
return []ent.Edge{
		edge.From("booking", Booking.Type).
			Ref("payments").
			Unique().
			Required(),
	}
}
