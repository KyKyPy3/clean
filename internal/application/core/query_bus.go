//nolint:dupl // it different structure
package core

import (
	"context"
	"errors"
	"fmt"
)

var (
	// ErrUnexpectedQuery represents the "unexpected query" error.
	ErrUnexpectedQuery = errors.New("unexpected query")

	// ErrQueryHandlerNotFound represents the "query handler not found" error.
	ErrQueryHandlerNotFound = errors.New("query handler not found")
)

// QueryType represents the type of a query.
type QueryType string

// Query represents the interface for queries.
type Query interface {
	Type() QueryType
}

// QueryHandler represents the interface for a query handler.
type QueryHandler interface {
	Handle(context.Context, Query) (any, error)
}

// QueryBus represents the query bus.
type QueryBus struct {
	handlers map[QueryType]QueryHandler
}

// NewQueryBus creates a new instance of QueryBus.
func NewQueryBus() *QueryBus {
	return &QueryBus{
		handlers: make(map[QueryType]QueryHandler),
	}
}

// Ask sends a query for processing.
func (c QueryBus) Ask(ctx context.Context, query Query) (any, error) {
	handler, ok := c.handlers[query.Type()]
	if !ok {
		return nil, fmt.Errorf("%s: %w", query.Type(), ErrQueryHandlerNotFound)
	}

	return handler.Handle(ctx, query)
}

// Register registers a query handler.
func (c QueryBus) Register(queryType QueryType, handler QueryHandler) {
	c.handlers[queryType] = handler
}
