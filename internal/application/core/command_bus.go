package core

import (
	"context"
	"errors"
	"fmt"
)

var (
	// ErrUnexpectedCommand represents the "unexpected command" error.
	ErrUnexpectedCommand = errors.New("unexpected command")

	// ErrCommandHandlerNotFound represents the "command handler not found" error.
	ErrCommandHandlerNotFound = errors.New("command handler not found")
)

// CommandType represents the type of a command.
type CommandType string

// Command represents the interface for commands.
type Command interface {
	Type() CommandType
}

// CommandHandler represents the interface for a command handler.
type CommandHandler interface {
	Handle(context.Context, Command) (any, error)
}

// CommandBus represents the command bus.
type CommandBus struct {
	handlers map[CommandType]CommandHandler
}

// NewCommandBus creates a new instance of CommandBus.
func NewCommandBus() *CommandBus {
	return &CommandBus{
		handlers: make(map[CommandType]CommandHandler),
	}
}

// Dispatch sends a command for processing.
func (c CommandBus) Dispatch(ctx context.Context, command Command) (any, error) {
	handler, ok := c.handlers[command.Type()]
	if !ok {
		return nil, fmt.Errorf("%s: %w", command.Type(), ErrCommandHandlerNotFound)
	}

	return handler.Handle(ctx, command)
}

// Register registers a command handler.
func (c CommandBus) Register(commandType CommandType, handler CommandHandler) {
	c.handlers[commandType] = handler
}
