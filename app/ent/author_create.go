// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"
	"time"

	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/patyhank/ptd/app/ent/author"
	"github.com/patyhank/ptd/app/ent/message"
)

// AuthorCreate is the builder for creating a Author entity.
type AuthorCreate struct {
	config
	mutation *AuthorMutation
	hooks    []Hook
}

// SetAuthorID sets the "author_id" field.
func (ac *AuthorCreate) SetAuthorID(s string) *AuthorCreate {
	ac.mutation.SetAuthorID(s)
	return ac
}

// SetLastSeen sets the "last_seen" field.
func (ac *AuthorCreate) SetLastSeen(t time.Time) *AuthorCreate {
	ac.mutation.SetLastSeen(t)
	return ac
}

// SetNillableLastSeen sets the "last_seen" field if the given value is not nil.
func (ac *AuthorCreate) SetNillableLastSeen(t *time.Time) *AuthorCreate {
	if t != nil {
		ac.SetLastSeen(*t)
	}
	return ac
}

// AddMessageIDs adds the "messages" edge to the Message entity by IDs.
func (ac *AuthorCreate) AddMessageIDs(ids ...int) *AuthorCreate {
	ac.mutation.AddMessageIDs(ids...)
	return ac
}

// AddMessages adds the "messages" edges to the Message entity.
func (ac *AuthorCreate) AddMessages(m ...*Message) *AuthorCreate {
	ids := make([]int, len(m))
	for i := range m {
		ids[i] = m[i].ID
	}
	return ac.AddMessageIDs(ids...)
}

// Mutation returns the AuthorMutation object of the builder.
func (ac *AuthorCreate) Mutation() *AuthorMutation {
	return ac.mutation
}

// Save creates the Author in the database.
func (ac *AuthorCreate) Save(ctx context.Context) (*Author, error) {
	ac.defaults()
	return withHooks(ctx, ac.sqlSave, ac.mutation, ac.hooks)
}

// SaveX calls Save and panics if Save returns an error.
func (ac *AuthorCreate) SaveX(ctx context.Context) *Author {
	v, err := ac.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (ac *AuthorCreate) Exec(ctx context.Context) error {
	_, err := ac.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (ac *AuthorCreate) ExecX(ctx context.Context) {
	if err := ac.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (ac *AuthorCreate) defaults() {
	if _, ok := ac.mutation.LastSeen(); !ok {
		v := author.DefaultLastSeen()
		ac.mutation.SetLastSeen(v)
	}
}

// check runs all checks and user-defined validators on the builder.
func (ac *AuthorCreate) check() error {
	if _, ok := ac.mutation.AuthorID(); !ok {
		return &ValidationError{Name: "author_id", err: errors.New(`ent: missing required field "Author.author_id"`)}
	}
	if _, ok := ac.mutation.LastSeen(); !ok {
		return &ValidationError{Name: "last_seen", err: errors.New(`ent: missing required field "Author.last_seen"`)}
	}
	return nil
}

func (ac *AuthorCreate) sqlSave(ctx context.Context) (*Author, error) {
	if err := ac.check(); err != nil {
		return nil, err
	}
	_node, _spec := ac.createSpec()
	if err := sqlgraph.CreateNode(ctx, ac.driver, _spec); err != nil {
		if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	id := _spec.ID.Value.(int64)
	_node.ID = int(id)
	ac.mutation.id = &_node.ID
	ac.mutation.done = true
	return _node, nil
}

func (ac *AuthorCreate) createSpec() (*Author, *sqlgraph.CreateSpec) {
	var (
		_node = &Author{config: ac.config}
		_spec = sqlgraph.NewCreateSpec(author.Table, sqlgraph.NewFieldSpec(author.FieldID, field.TypeInt))
	)
	if value, ok := ac.mutation.AuthorID(); ok {
		_spec.SetField(author.FieldAuthorID, field.TypeString, value)
		_node.AuthorID = value
	}
	if value, ok := ac.mutation.LastSeen(); ok {
		_spec.SetField(author.FieldLastSeen, field.TypeTime, value)
		_node.LastSeen = value
	}
	if nodes := ac.mutation.MessagesIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: true,
			Table:   author.MessagesTable,
			Columns: []string{author.MessagesColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(message.FieldID, field.TypeInt),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges = append(_spec.Edges, edge)
	}
	return _node, _spec
}

// AuthorCreateBulk is the builder for creating many Author entities in bulk.
type AuthorCreateBulk struct {
	config
	err      error
	builders []*AuthorCreate
}

// Save creates the Author entities in the database.
func (acb *AuthorCreateBulk) Save(ctx context.Context) ([]*Author, error) {
	if acb.err != nil {
		return nil, acb.err
	}
	specs := make([]*sqlgraph.CreateSpec, len(acb.builders))
	nodes := make([]*Author, len(acb.builders))
	mutators := make([]Mutator, len(acb.builders))
	for i := range acb.builders {
		func(i int, root context.Context) {
			builder := acb.builders[i]
			builder.defaults()
			var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
				mutation, ok := m.(*AuthorMutation)
				if !ok {
					return nil, fmt.Errorf("unexpected mutation type %T", m)
				}
				if err := builder.check(); err != nil {
					return nil, err
				}
				builder.mutation = mutation
				var err error
				nodes[i], specs[i] = builder.createSpec()
				if i < len(mutators)-1 {
					_, err = mutators[i+1].Mutate(root, acb.builders[i+1].mutation)
				} else {
					spec := &sqlgraph.BatchCreateSpec{Nodes: specs}
					// Invoke the actual operation on the latest mutation in the chain.
					if err = sqlgraph.BatchCreate(ctx, acb.driver, spec); err != nil {
						if sqlgraph.IsConstraintError(err) {
							err = &ConstraintError{msg: err.Error(), wrap: err}
						}
					}
				}
				if err != nil {
					return nil, err
				}
				mutation.id = &nodes[i].ID
				if specs[i].ID.Value != nil {
					id := specs[i].ID.Value.(int64)
					nodes[i].ID = int(id)
				}
				mutation.done = true
				return nodes[i], nil
			})
			for i := len(builder.hooks) - 1; i >= 0; i-- {
				mut = builder.hooks[i](mut)
			}
			mutators[i] = mut
		}(i, ctx)
	}
	if len(mutators) > 0 {
		if _, err := mutators[0].Mutate(ctx, acb.builders[0].mutation); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

// SaveX is like Save, but panics if an error occurs.
func (acb *AuthorCreateBulk) SaveX(ctx context.Context) []*Author {
	v, err := acb.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (acb *AuthorCreateBulk) Exec(ctx context.Context) error {
	_, err := acb.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (acb *AuthorCreateBulk) ExecX(ctx context.Context) {
	if err := acb.Exec(ctx); err != nil {
		panic(err)
	}
}
