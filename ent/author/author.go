// Code generated by ent, DO NOT EDIT.

package author

import (
	"time"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
)

const (
	// Label holds the string label denoting the author type in the database.
	Label = "author"
	// FieldID holds the string denoting the id field in the database.
	FieldID = "id"
	// FieldAuthorID holds the string denoting the author_id field in the database.
	FieldAuthorID = "author_id"
	// FieldLastSeen holds the string denoting the last_seen field in the database.
	FieldLastSeen = "last_seen"
	// EdgeMessages holds the string denoting the messages edge name in mutations.
	EdgeMessages = "messages"
	// Table holds the table name of the author in the database.
	Table = "authors"
	// MessagesTable is the table that holds the messages relation/edge.
	MessagesTable = "messages"
	// MessagesInverseTable is the table name for the Message entity.
	// It exists in this package in order to avoid circular dependency with the "message" package.
	MessagesInverseTable = "messages"
	// MessagesColumn is the table column denoting the messages relation/edge.
	MessagesColumn = "message_author"
)

// Columns holds all SQL columns for author fields.
var Columns = []string{
	FieldID,
	FieldAuthorID,
	FieldLastSeen,
}

// ValidColumn reports if the column name is valid (part of the table columns).
func ValidColumn(column string) bool {
	for i := range Columns {
		if column == Columns[i] {
			return true
		}
	}
	return false
}

var (
	// DefaultLastSeen holds the default value on creation for the "last_seen" field.
	DefaultLastSeen func() time.Time
)

// OrderOption defines the ordering options for the Author queries.
type OrderOption func(*sql.Selector)

// ByID orders the results by the id field.
func ByID(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldID, opts...).ToFunc()
}

// ByAuthorID orders the results by the author_id field.
func ByAuthorID(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldAuthorID, opts...).ToFunc()
}

// ByLastSeen orders the results by the last_seen field.
func ByLastSeen(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldLastSeen, opts...).ToFunc()
}

// ByMessagesCount orders the results by messages count.
func ByMessagesCount(opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborsCount(s, newMessagesStep(), opts...)
	}
}

// ByMessages orders the results by messages terms.
func ByMessages(term sql.OrderTerm, terms ...sql.OrderTerm) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newMessagesStep(), append([]sql.OrderTerm{term}, terms...)...)
	}
}
func newMessagesStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(MessagesInverseTable, FieldID),
		sqlgraph.Edge(sqlgraph.O2M, true, MessagesTable, MessagesColumn),
	)
}
