package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

// Airport holds the schema definition for the Airport entity.
type Airport struct {
	ent.Schema
}

// Fields of the Airport.
func (Airport) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(uuid.New),

		field.String("code").
			Unique().
			MaxLen(3),

		field.String("name"),
		field.String("city"),
		field.String("country"),

		field.Time("created_at").
			Default(time.Now),
	}
}

// Edges of the Airport.
func (Airport) Edges() []ent.Edge {
		return []ent.Edge{
		edge.To("departures", Flight.Type),
		edge.To("arrivals", Flight.Type),
	}
}
