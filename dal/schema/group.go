package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

type Group struct {
	ent.Schema
}

func (Group) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").NotEmpty().Unique(),
	}
}

func (Group) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("cards", Card.Type).Ref("group"),
		edge.To("types", Typ3.Type),
		edge.To("supertype", Supertype.Type).Unique(),
	}
}

func (Group) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("name").Unique(),
	}
}
