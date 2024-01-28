package core

import "github.com/KyKyPy3/clean/pkg/mediator"

type AggregateRoot interface {
	AddEvent(event mediator.Event)
	Events() []mediator.Event
}

type BaseAggregateRoot struct {
	events []mediator.Event
}

func (a *BaseAggregateRoot) AddEvent(event mediator.Event) {
	a.events = append(a.events, event)
}

func (a *BaseAggregateRoot) Events() []mediator.Event {
	events := a.events
	a.events = nil

	return events
}

var _ AggregateRoot = (*BaseAggregateRoot)(nil)
