package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

type Typ3 struct {
	ent.Schema
}

func (Typ3) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").NotEmpty().Unique(),
	}
}

func (Typ3) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("card", Card.Type).Ref("types"),
		edge.From("group", Group.Type).Ref("types"),
	}
}

func (Typ3) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("name").Unique(),
	}
}
