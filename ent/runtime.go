// Code generated by ent, DO NOT EDIT.

package ent

import (
	"time"

	snowflake "github.com/disgoorg/snowflake/v2"
	"github.com/patyhank/ptd/ent/author"
	"github.com/patyhank/ptd/ent/message"
	"github.com/patyhank/ptd/ent/postinfo"
	"github.com/patyhank/ptd/ent/schema"
)

// The init function reads all schema descriptors with runtime code
// (default values, validators, hooks and policies) and stitches it
// to their package variables.
func init() {
	authorFields := schema.Author{}.Fields()
	_ = authorFields
	// authorDescLastSeen is the schema descriptor for last_seen field.
	authorDescLastSeen := authorFields[1].Descriptor()
	// author.DefaultLastSeen holds the default value on creation for the last_seen field.
	author.DefaultLastSeen = authorDescLastSeen.Default.(func() time.Time)
	messageFields := schema.Message{}.Fields()
	_ = messageFields
	// messageDescMessageID is the schema descriptor for message_id field.
	messageDescMessageID := messageFields[4].Descriptor()
	// message.DefaultMessageID holds the default value on creation for the message_id field.
	message.DefaultMessageID = snowflake.ID(messageDescMessageID.Default.(uint64))
	// message.MessageIDValidator is a validator for the "message_id" field. It is called by the builders before save.
	message.MessageIDValidator = messageDescMessageID.Validators[0].(func(uint64) error)
	postinfoFields := schema.PostInfo{}.Fields()
	_ = postinfoFields
	// postinfoDescLastUpdated is the schema descriptor for last_updated field.
	postinfoDescLastUpdated := postinfoFields[1].Descriptor()
	// postinfo.DefaultLastUpdated holds the default value on creation for the last_updated field.
	postinfo.DefaultLastUpdated = postinfoDescLastUpdated.Default.(func() time.Time)
	// postinfoDescCurrentViewing is the schema descriptor for current_viewing field.
	postinfoDescCurrentViewing := postinfoFields[2].Descriptor()
	// postinfo.DefaultCurrentViewing holds the default value on creation for the current_viewing field.
	postinfo.DefaultCurrentViewing = postinfoDescCurrentViewing.Default.(bool)
	// postinfoDescShouldArchived is the schema descriptor for should_archived field.
	postinfoDescShouldArchived := postinfoFields[4].Descriptor()
	// postinfo.DefaultShouldArchived holds the default value on creation for the should_archived field.
	postinfo.DefaultShouldArchived = postinfoDescShouldArchived.Default.(bool)
	// postinfoDescChannelID is the schema descriptor for channel_id field.
	postinfoDescChannelID := postinfoFields[9].Descriptor()
	// postinfo.DefaultChannelID holds the default value on creation for the channel_id field.
	postinfo.DefaultChannelID = snowflake.ID(postinfoDescChannelID.Default.(uint64))
	// postinfo.ChannelIDValidator is a validator for the "channel_id" field. It is called by the builders before save.
	postinfo.ChannelIDValidator = postinfoDescChannelID.Validators[0].(func(uint64) error)
}
