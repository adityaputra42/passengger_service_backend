package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

// Aircraft holds the schema definition for the Aircraft entity.
type Aircraft struct {
	ent.Schema
}

// Fields of the Aircraft.
func (Aircraft) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(uuid.New),

		field.String("model"),
		field.String("registration_no").
			Unique(),

		field.Int("total_seats"),

		field.Time("created_at").
			Default(time.Now),
	}
}

// Edges of the Aircraft.
func (Aircraft) Edges() []ent.Edge {
		return []ent.Edge{
		edge.To("flights", Flight.Type),
		edge.To("seats", Seats.Type),

	}
}
