package common

type Event interface {
	Kind() string
}

type Entity interface {
	ID() ID
	SetID(id ID)
	Events() []Event
}

type BaseEntity struct {
	id     ID
	events []Event
}

func (e *BaseEntity) ID() ID {
	return e.id
}

func (e *BaseEntity) SetID(id ID) {
	e.id = id
}

func (e *BaseEntity) Events() []Event {
	return e.events
}

var _ Entity = (*BaseEntity)(nil)
