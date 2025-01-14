// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"
	"time"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	snowflake "github.com/disgoorg/snowflake/v2"
	"github.com/patyhank/ptd/app/ent/author"
	"github.com/patyhank/ptd/app/ent/message"
	"github.com/patyhank/ptd/app/ent/postinfo"
	"github.com/patyhank/ptd/app/ent/predicate"
)

// MessageUpdate is the builder for updating Message entities.
type MessageUpdate struct {
	config
	hooks    []Hook
	mutation *MessageMutation
}

// Where appends a list predicates to the MessageUpdate builder.
func (mu *MessageUpdate) Where(ps ...predicate.Message) *MessageUpdate {
	mu.mutation.Where(ps...)
	return mu
}

// SetHash sets the "hash" field.
func (mu *MessageUpdate) SetHash(s string) *MessageUpdate {
	mu.mutation.SetHash(s)
	return mu
}

// SetNillableHash sets the "hash" field if the given value is not nil.
func (mu *MessageUpdate) SetNillableHash(s *string) *MessageUpdate {
	if s != nil {
		mu.SetHash(*s)
	}
	return mu
}

// SetContent sets the "content" field.
func (mu *MessageUpdate) SetContent(s string) *MessageUpdate {
	mu.mutation.SetContent(s)
	return mu
}

// SetNillableContent sets the "content" field if the given value is not nil.
func (mu *MessageUpdate) SetNillableContent(s *string) *MessageUpdate {
	if s != nil {
		mu.SetContent(*s)
	}
	return mu
}

// SetRawContent sets the "raw_content" field.
func (mu *MessageUpdate) SetRawContent(s string) *MessageUpdate {
	mu.mutation.SetRawContent(s)
	return mu
}

// SetNillableRawContent sets the "raw_content" field if the given value is not nil.
func (mu *MessageUpdate) SetNillableRawContent(s *string) *MessageUpdate {
	if s != nil {
		mu.SetRawContent(*s)
	}
	return mu
}

// SetCreatedAt sets the "created_at" field.
func (mu *MessageUpdate) SetCreatedAt(t time.Time) *MessageUpdate {
	mu.mutation.SetCreatedAt(t)
	return mu
}

// SetNillableCreatedAt sets the "created_at" field if the given value is not nil.
func (mu *MessageUpdate) SetNillableCreatedAt(t *time.Time) *MessageUpdate {
	if t != nil {
		mu.SetCreatedAt(*t)
	}
	return mu
}

// SetMessageID sets the "message_id" field.
func (mu *MessageUpdate) SetMessageID(s snowflake.ID) *MessageUpdate {
	mu.mutation.ResetMessageID()
	mu.mutation.SetMessageID(s)
	return mu
}

// SetNillableMessageID sets the "message_id" field if the given value is not nil.
func (mu *MessageUpdate) SetNillableMessageID(s *snowflake.ID) *MessageUpdate {
	if s != nil {
		mu.SetMessageID(*s)
	}
	return mu
}

// AddMessageID adds s to the "message_id" field.
func (mu *MessageUpdate) AddMessageID(s snowflake.ID) *MessageUpdate {
	mu.mutation.AddMessageID(s)
	return mu
}

// SetAuthorID sets the "author" edge to the Author entity by ID.
func (mu *MessageUpdate) SetAuthorID(id int) *MessageUpdate {
	mu.mutation.SetAuthorID(id)
	return mu
}

// SetAuthor sets the "author" edge to the Author entity.
func (mu *MessageUpdate) SetAuthor(a *Author) *MessageUpdate {
	return mu.SetAuthorID(a.ID)
}

// SetParentPostID sets the "parent_post" edge to the PostInfo entity by ID.
func (mu *MessageUpdate) SetParentPostID(id int) *MessageUpdate {
	mu.mutation.SetParentPostID(id)
	return mu
}

// SetParentPost sets the "parent_post" edge to the PostInfo entity.
func (mu *MessageUpdate) SetParentPost(p *PostInfo) *MessageUpdate {
	return mu.SetParentPostID(p.ID)
}

// Mutation returns the MessageMutation object of the builder.
func (mu *MessageUpdate) Mutation() *MessageMutation {
	return mu.mutation
}

// ClearAuthor clears the "author" edge to the Author entity.
func (mu *MessageUpdate) ClearAuthor() *MessageUpdate {
	mu.mutation.ClearAuthor()
	return mu
}

// ClearParentPost clears the "parent_post" edge to the PostInfo entity.
func (mu *MessageUpdate) ClearParentPost() *MessageUpdate {
	mu.mutation.ClearParentPost()
	return mu
}

// Save executes the query and returns the number of nodes affected by the update operation.
func (mu *MessageUpdate) Save(ctx context.Context) (int, error) {
	return withHooks(ctx, mu.sqlSave, mu.mutation, mu.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (mu *MessageUpdate) SaveX(ctx context.Context) int {
	affected, err := mu.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (mu *MessageUpdate) Exec(ctx context.Context) error {
	_, err := mu.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (mu *MessageUpdate) ExecX(ctx context.Context) {
	if err := mu.Exec(ctx); err != nil {
		panic(err)
	}
}

// check runs all checks and user-defined validators on the builder.
func (mu *MessageUpdate) check() error {
	if v, ok := mu.mutation.MessageID(); ok {
		if err := message.MessageIDValidator(uint64(v)); err != nil {
			return &ValidationError{Name: "message_id", err: fmt.Errorf(`ent: validator failed for field "Message.message_id": %w`, err)}
		}
	}
	if mu.mutation.AuthorCleared() && len(mu.mutation.AuthorIDs()) > 0 {
		return errors.New(`ent: clearing a required unique edge "Message.author"`)
	}
	if mu.mutation.ParentPostCleared() && len(mu.mutation.ParentPostIDs()) > 0 {
		return errors.New(`ent: clearing a required unique edge "Message.parent_post"`)
	}
	return nil
}

func (mu *MessageUpdate) sqlSave(ctx context.Context) (n int, err error) {
	if err := mu.check(); err != nil {
		return n, err
	}
	_spec := sqlgraph.NewUpdateSpec(message.Table, message.Columns, sqlgraph.NewFieldSpec(message.FieldID, field.TypeInt))
	if ps := mu.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := mu.mutation.Hash(); ok {
		_spec.SetField(message.FieldHash, field.TypeString, value)
	}
	if value, ok := mu.mutation.Content(); ok {
		_spec.SetField(message.FieldContent, field.TypeString, value)
	}
	if value, ok := mu.mutation.RawContent(); ok {
		_spec.SetField(message.FieldRawContent, field.TypeString, value)
	}
	if value, ok := mu.mutation.CreatedAt(); ok {
		_spec.SetField(message.FieldCreatedAt, field.TypeTime, value)
	}
	if value, ok := mu.mutation.MessageID(); ok {
		_spec.SetField(message.FieldMessageID, field.TypeUint64, value)
	}
	if value, ok := mu.mutation.AddedMessageID(); ok {
		_spec.AddField(message.FieldMessageID, field.TypeUint64, value)
	}
	if mu.mutation.AuthorCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   message.AuthorTable,
			Columns: []string{message.AuthorColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(author.FieldID, field.TypeInt),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := mu.mutation.AuthorIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   message.AuthorTable,
			Columns: []string{message.AuthorColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(author.FieldID, field.TypeInt),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if mu.mutation.ParentPostCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   message.ParentPostTable,
			Columns: []string{message.ParentPostColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(postinfo.FieldID, field.TypeInt),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := mu.mutation.ParentPostIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   message.ParentPostTable,
			Columns: []string{message.ParentPostColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(postinfo.FieldID, field.TypeInt),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if n, err = sqlgraph.UpdateNodes(ctx, mu.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{message.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return 0, err
	}
	mu.mutation.done = true
	return n, nil
}

// MessageUpdateOne is the builder for updating a single Message entity.
type MessageUpdateOne struct {
	config
	fields   []string
	hooks    []Hook
	mutation *MessageMutation
}

// SetHash sets the "hash" field.
func (muo *MessageUpdateOne) SetHash(s string) *MessageUpdateOne {
	muo.mutation.SetHash(s)
	return muo
}

// SetNillableHash sets the "hash" field if the given value is not nil.
func (muo *MessageUpdateOne) SetNillableHash(s *string) *MessageUpdateOne {
	if s != nil {
		muo.SetHash(*s)
	}
	return muo
}

// SetContent sets the "content" field.
func (muo *MessageUpdateOne) SetContent(s string) *MessageUpdateOne {
	muo.mutation.SetContent(s)
	return muo
}

// SetNillableContent sets the "content" field if the given value is not nil.
func (muo *MessageUpdateOne) SetNillableContent(s *string) *MessageUpdateOne {
	if s != nil {
		muo.SetContent(*s)
	}
	return muo
}

// SetRawContent sets the "raw_content" field.
func (muo *MessageUpdateOne) SetRawContent(s string) *MessageUpdateOne {
	muo.mutation.SetRawContent(s)
	return muo
}

// SetNillableRawContent sets the "raw_content" field if the given value is not nil.
func (muo *MessageUpdateOne) SetNillableRawContent(s *string) *MessageUpdateOne {
	if s != nil {
		muo.SetRawContent(*s)
	}
	return muo
}

// SetCreatedAt sets the "created_at" field.
func (muo *MessageUpdateOne) SetCreatedAt(t time.Time) *MessageUpdateOne {
	muo.mutation.SetCreatedAt(t)
	return muo
}

// SetNillableCreatedAt sets the "created_at" field if the given value is not nil.
func (muo *MessageUpdateOne) SetNillableCreatedAt(t *time.Time) *MessageUpdateOne {
	if t != nil {
		muo.SetCreatedAt(*t)
	}
	return muo
}

// SetMessageID sets the "message_id" field.
func (muo *MessageUpdateOne) SetMessageID(s snowflake.ID) *MessageUpdateOne {
	muo.mutation.ResetMessageID()
	muo.mutation.SetMessageID(s)
	return muo
}

// SetNillableMessageID sets the "message_id" field if the given value is not nil.
func (muo *MessageUpdateOne) SetNillableMessageID(s *snowflake.ID) *MessageUpdateOne {
	if s != nil {
		muo.SetMessageID(*s)
	}
	return muo
}

// AddMessageID adds s to the "message_id" field.
func (muo *MessageUpdateOne) AddMessageID(s snowflake.ID) *MessageUpdateOne {
	muo.mutation.AddMessageID(s)
	return muo
}

// SetAuthorID sets the "author" edge to the Author entity by ID.
func (muo *MessageUpdateOne) SetAuthorID(id int) *MessageUpdateOne {
	muo.mutation.SetAuthorID(id)
	return muo
}

// SetAuthor sets the "author" edge to the Author entity.
func (muo *MessageUpdateOne) SetAuthor(a *Author) *MessageUpdateOne {
	return muo.SetAuthorID(a.ID)
}

// SetParentPostID sets the "parent_post" edge to the PostInfo entity by ID.
func (muo *MessageUpdateOne) SetParentPostID(id int) *MessageUpdateOne {
	muo.mutation.SetParentPostID(id)
	return muo
}

// SetParentPost sets the "parent_post" edge to the PostInfo entity.
func (muo *MessageUpdateOne) SetParentPost(p *PostInfo) *MessageUpdateOne {
	return muo.SetParentPostID(p.ID)
}

// Mutation returns the MessageMutation object of the builder.
func (muo *MessageUpdateOne) Mutation() *MessageMutation {
	return muo.mutation
}

// ClearAuthor clears the "author" edge to the Author entity.
func (muo *MessageUpdateOne) ClearAuthor() *MessageUpdateOne {
	muo.mutation.ClearAuthor()
	return muo
}

// ClearParentPost clears the "parent_post" edge to the PostInfo entity.
func (muo *MessageUpdateOne) ClearParentPost() *MessageUpdateOne {
	muo.mutation.ClearParentPost()
	return muo
}

// Where appends a list predicates to the MessageUpdate builder.
func (muo *MessageUpdateOne) Where(ps ...predicate.Message) *MessageUpdateOne {
	muo.mutation.Where(ps...)
	return muo
}

// Select allows selecting one or more fields (columns) of the returned entity.
// The default is selecting all fields defined in the entity schema.
func (muo *MessageUpdateOne) Select(field string, fields ...string) *MessageUpdateOne {
	muo.fields = append([]string{field}, fields...)
	return muo
}

// Save executes the query and returns the updated Message entity.
func (muo *MessageUpdateOne) Save(ctx context.Context) (*Message, error) {
	return withHooks(ctx, muo.sqlSave, muo.mutation, muo.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (muo *MessageUpdateOne) SaveX(ctx context.Context) *Message {
	node, err := muo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// Exec executes the query on the entity.
func (muo *MessageUpdateOne) Exec(ctx context.Context) error {
	_, err := muo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (muo *MessageUpdateOne) ExecX(ctx context.Context) {
	if err := muo.Exec(ctx); err != nil {
		panic(err)
	}
}

// check runs all checks and user-defined validators on the builder.
func (muo *MessageUpdateOne) check() error {
	if v, ok := muo.mutation.MessageID(); ok {
		if err := message.MessageIDValidator(uint64(v)); err != nil {
			return &ValidationError{Name: "message_id", err: fmt.Errorf(`ent: validator failed for field "Message.message_id": %w`, err)}
		}
	}
	if muo.mutation.AuthorCleared() && len(muo.mutation.AuthorIDs()) > 0 {
		return errors.New(`ent: clearing a required unique edge "Message.author"`)
	}
	if muo.mutation.ParentPostCleared() && len(muo.mutation.ParentPostIDs()) > 0 {
		return errors.New(`ent: clearing a required unique edge "Message.parent_post"`)
	}
	return nil
}

func (muo *MessageUpdateOne) sqlSave(ctx context.Context) (_node *Message, err error) {
	if err := muo.check(); err != nil {
		return _node, err
	}
	_spec := sqlgraph.NewUpdateSpec(message.Table, message.Columns, sqlgraph.NewFieldSpec(message.FieldID, field.TypeInt))
	id, ok := muo.mutation.ID()
	if !ok {
		return nil, &ValidationError{Name: "id", err: errors.New(`ent: missing "Message.id" for update`)}
	}
	_spec.Node.ID.Value = id
	if fields := muo.fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, message.FieldID)
		for _, f := range fields {
			if !message.ValidColumn(f) {
				return nil, &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
			}
			if f != message.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, f)
			}
		}
	}
	if ps := muo.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := muo.mutation.Hash(); ok {
		_spec.SetField(message.FieldHash, field.TypeString, value)
	}
	if value, ok := muo.mutation.Content(); ok {
		_spec.SetField(message.FieldContent, field.TypeString, value)
	}
	if value, ok := muo.mutation.RawContent(); ok {
		_spec.SetField(message.FieldRawContent, field.TypeString, value)
	}
	if value, ok := muo.mutation.CreatedAt(); ok {
		_spec.SetField(message.FieldCreatedAt, field.TypeTime, value)
	}
	if value, ok := muo.mutation.MessageID(); ok {
		_spec.SetField(message.FieldMessageID, field.TypeUint64, value)
	}
	if value, ok := muo.mutation.AddedMessageID(); ok {
		_spec.AddField(message.FieldMessageID, field.TypeUint64, value)
	}
	if muo.mutation.AuthorCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   message.AuthorTable,
			Columns: []string{message.AuthorColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(author.FieldID, field.TypeInt),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := muo.mutation.AuthorIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   message.AuthorTable,
			Columns: []string{message.AuthorColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(author.FieldID, field.TypeInt),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if muo.mutation.ParentPostCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   message.ParentPostTable,
			Columns: []string{message.ParentPostColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(postinfo.FieldID, field.TypeInt),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := muo.mutation.ParentPostIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   message.ParentPostTable,
			Columns: []string{message.ParentPostColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(postinfo.FieldID, field.TypeInt),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	_node = &Message{config: muo.config}
	_spec.Assign = _node.assignValues
	_spec.ScanValues = _node.scanValues
	if err = sqlgraph.UpdateNode(ctx, muo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{message.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	muo.mutation.done = true
	return _node, nil
}
