// Code generated by ent, DO NOT EDIT.

package ent

import (
	"fmt"
	"strings"
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/sql"
	snowflake "github.com/disgoorg/snowflake/v2"
	"github.com/patyhank/ptd/app/ent/author"
	"github.com/patyhank/ptd/app/ent/message"
	"github.com/patyhank/ptd/app/ent/postinfo"
)

// Message is the model entity for the Message schema.
type Message struct {
	config `json:"-"`
	// ID of the ent.
	ID int `json:"id,omitempty"`
	// Hash holds the value of the "hash" field.
	Hash string `json:"hash,omitempty"`
	// Content holds the value of the "content" field.
	Content string `json:"content,omitempty"`
	// RawContent holds the value of the "raw_content" field.
	RawContent string `json:"raw_content,omitempty"`
	// CreatedAt holds the value of the "created_at" field.
	CreatedAt time.Time `json:"created_at,omitempty"`
	// MessageID holds the value of the "message_id" field.
	MessageID snowflake.ID `json:"message_id,omitempty"`
	// Edges holds the relations/edges for other nodes in the graph.
	// The values are being populated by the MessageQuery when eager-loading is set.
	Edges              MessageEdges `json:"edges"`
	message_author     *int
	post_info_messages *int
	selectValues       sql.SelectValues
}

// MessageEdges holds the relations/edges for other nodes in the graph.
type MessageEdges struct {
	// Author holds the value of the author edge.
	Author *Author `json:"author,omitempty"`
	// ParentPost holds the value of the parent_post edge.
	ParentPost *PostInfo `json:"parent_post,omitempty"`
	// loadedTypes holds the information for reporting if a
	// type was loaded (or requested) in eager-loading or not.
	loadedTypes [2]bool
}

// AuthorOrErr returns the Author value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e MessageEdges) AuthorOrErr() (*Author, error) {
	if e.Author != nil {
		return e.Author, nil
	} else if e.loadedTypes[0] {
		return nil, &NotFoundError{label: author.Label}
	}
	return nil, &NotLoadedError{edge: "author"}
}

// ParentPostOrErr returns the ParentPost value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e MessageEdges) ParentPostOrErr() (*PostInfo, error) {
	if e.ParentPost != nil {
		return e.ParentPost, nil
	} else if e.loadedTypes[1] {
		return nil, &NotFoundError{label: postinfo.Label}
	}
	return nil, &NotLoadedError{edge: "parent_post"}
}

// scanValues returns the types for scanning values from sql.Rows.
func (*Message) scanValues(columns []string) ([]any, error) {
	values := make([]any, len(columns))
	for i := range columns {
		switch columns[i] {
		case message.FieldID, message.FieldMessageID:
			values[i] = new(sql.NullInt64)
		case message.FieldHash, message.FieldContent, message.FieldRawContent:
			values[i] = new(sql.NullString)
		case message.FieldCreatedAt:
			values[i] = new(sql.NullTime)
		case message.ForeignKeys[0]: // message_author
			values[i] = new(sql.NullInt64)
		case message.ForeignKeys[1]: // post_info_messages
			values[i] = new(sql.NullInt64)
		default:
			values[i] = new(sql.UnknownType)
		}
	}
	return values, nil
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the Message fields.
func (m *Message) assignValues(columns []string, values []any) error {
	if m, n := len(values), len(columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	for i := range columns {
		switch columns[i] {
		case message.FieldID:
			value, ok := values[i].(*sql.NullInt64)
			if !ok {
				return fmt.Errorf("unexpected type %T for field id", value)
			}
			m.ID = int(value.Int64)
		case message.FieldHash:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field hash", values[i])
			} else if value.Valid {
				m.Hash = value.String
			}
		case message.FieldContent:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field content", values[i])
			} else if value.Valid {
				m.Content = value.String
			}
		case message.FieldRawContent:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field raw_content", values[i])
			} else if value.Valid {
				m.RawContent = value.String
			}
		case message.FieldCreatedAt:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field created_at", values[i])
			} else if value.Valid {
				m.CreatedAt = value.Time
			}
		case message.FieldMessageID:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for field message_id", values[i])
			} else if value.Valid {
				m.MessageID = snowflake.ID(value.Int64)
			}
		case message.ForeignKeys[0]:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for edge-field message_author", value)
			} else if value.Valid {
				m.message_author = new(int)
				*m.message_author = int(value.Int64)
			}
		case message.ForeignKeys[1]:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for edge-field post_info_messages", value)
			} else if value.Valid {
				m.post_info_messages = new(int)
				*m.post_info_messages = int(value.Int64)
			}
		default:
			m.selectValues.Set(columns[i], values[i])
		}
	}
	return nil
}

// Value returns the ent.Value that was dynamically selected and assigned to the Message.
// This includes values selected through modifiers, order, etc.
func (m *Message) Value(name string) (ent.Value, error) {
	return m.selectValues.Get(name)
}

// QueryAuthor queries the "author" edge of the Message entity.
func (m *Message) QueryAuthor() *AuthorQuery {
	return NewMessageClient(m.config).QueryAuthor(m)
}

// QueryParentPost queries the "parent_post" edge of the Message entity.
func (m *Message) QueryParentPost() *PostInfoQuery {
	return NewMessageClient(m.config).QueryParentPost(m)
}

// Update returns a builder for updating this Message.
// Note that you need to call Message.Unwrap() before calling this method if this Message
// was returned from a transaction, and the transaction was committed or rolled back.
func (m *Message) Update() *MessageUpdateOne {
	return NewMessageClient(m.config).UpdateOne(m)
}

// Unwrap unwraps the Message entity that was returned from a transaction after it was closed,
// so that all future queries will be executed through the driver which created the transaction.
func (m *Message) Unwrap() *Message {
	_tx, ok := m.config.driver.(*txDriver)
	if !ok {
		panic("ent: Message is not a transactional entity")
	}
	m.config.driver = _tx.drv
	return m
}

// String implements the fmt.Stringer.
func (m *Message) String() string {
	var builder strings.Builder
	builder.WriteString("Message(")
	builder.WriteString(fmt.Sprintf("id=%v, ", m.ID))
	builder.WriteString("hash=")
	builder.WriteString(m.Hash)
	builder.WriteString(", ")
	builder.WriteString("content=")
	builder.WriteString(m.Content)
	builder.WriteString(", ")
	builder.WriteString("raw_content=")
	builder.WriteString(m.RawContent)
	builder.WriteString(", ")
	builder.WriteString("created_at=")
	builder.WriteString(m.CreatedAt.Format(time.ANSIC))
	builder.WriteString(", ")
	builder.WriteString("message_id=")
	builder.WriteString(fmt.Sprintf("%v", m.MessageID))
	builder.WriteByte(')')
	return builder.String()
}

// Messages is a parsable slice of Message.
type Messages []*Message
