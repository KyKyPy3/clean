package mediator

import (
	"context"

	"github.com/KyKyPy3/clean/pkg/logger"
)

type Event interface {
	Kind() string
}

type Handler func(context.Context, Event) error

type Mediator struct {
	handlers map[string][]Handler
	logger   logger.Logger
}

func New(logger logger.Logger) *Mediator {
	return &Mediator{
		handlers: map[string][]Handler{},
		logger:   logger,
	}
}

func (m *Mediator) Subscribe(kind string, handler Handler) {
	handlers := m.handlers[kind]
	m.handlers[kind] = append(handlers, handler)
}

func (m *Mediator) Publish(ctx context.Context, events ...Event) error {
	for _, event := range events {
		handlers := m.handlers[event.Kind()]
		for _, handler := range handlers {
			if err := handler(ctx, event); err != nil {
				m.logger.Errorf("can't publish event %v, err: %v", event, err)
			}
		}
	}

	return nil
}
