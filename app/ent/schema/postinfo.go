package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/disgoorg/snowflake/v2"
	"time"
)

// PostInfo holds the schema definition for the PostInfo entity.
type PostInfo struct {
	ent.Schema
}

// Fields of the PostInfo.
func (PostInfo) Fields() []ent.Field {
	return []ent.Field{
		field.String("title").Optional().Unique(),
		field.Time("last_updated").Default(time.Now),
		field.Bool("current_viewing").Default(true),
		field.Strings("search_keywords").Optional(),
		field.Time("force_view_expire").Optional(),
		field.Bool("should_archived").Default(false),
		field.String("aid").Optional(),
		field.String("url").Optional(),
		field.Text("post_content").Optional(),
		field.JSON("content_messages", []snowflake.ID{}).Optional(),
		field.Uint64("channel_id").GoType(snowflake.ID(0)).Max(18446744073709551615).Default(0), // set max size to prevent EntGo constraint error
	}
}

// Edges of the PostInfo.
func (PostInfo) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("messages", Message.Type),
	}
}
