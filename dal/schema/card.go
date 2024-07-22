package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/Farmer-Pete/HokeyPoke/data"
)

type Card struct {
	ent.Schema
}

func (Card) Fields() []ent.Field {
	return []ent.Field{
		field.String("ptcg_id").NotEmpty().Unique(),
		field.String("name").NotEmpty(),
		field.JSON("metadata", data.CardMetadata{}),
	}
}

func (Card) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("supertype", Supertype.Type).Unique(),
		edge.To("types", Typ3.Type),
		edge.To("group", Group.Type).Unique(),
	}
}

func (Card) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("ptcg_id").Unique(),
	}
}
