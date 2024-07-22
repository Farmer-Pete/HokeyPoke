package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

type Supertype struct {
	ent.Schema
}

func (Supertype) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").NotEmpty().Unique(),
	}
}

func (Supertype) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("card", Card.Type).Ref("supertype"),
		edge.From("group", Group.Type).Ref("supertype"),
	}
}

func (Supertype) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("name").Unique(),
	}
}
