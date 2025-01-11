package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/disgoorg/snowflake/v2"
)

// Message holds the schema definition for the Message entity.
type Message struct {
	ent.Schema
}

// Fields of the Message.
func (Message) Fields() []ent.Field {
	return []ent.Field{
		field.String("hash"),
		field.Text("content"),
		field.Text("raw_content"),
		field.Time("created_at"),
		field.Uint64("message_id").GoType(snowflake.ID(0)).Max(18446744073709551615).Default(0), // set max size to prevent EntGo constraint error
	}
}

// Edges of the Message.
func (Message) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("author", Author.Type).Unique().Required(),
		edge.From("parent_post", PostInfo.Type).Ref("messages").Unique().Required(),
	}
}
