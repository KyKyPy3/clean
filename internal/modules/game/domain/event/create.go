package event

const GameCreated = "GameCreated"

type GameCreatedEvent struct {
	ID   string
	Name string
}

func (e GameCreatedEvent) Kind() string {
	return GameCreated
}
