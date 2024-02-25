package entity

import (
	"fmt"
	"time"

	"github.com/KyKyPy3/clean/internal/domain/common"
	"github.com/KyKyPy3/clean/internal/domain/core"
	"github.com/KyKyPy3/clean/internal/modules/game/domain/event"
	"github.com/KyKyPy3/clean/pkg/mediator"
)

// Game struct.
type Game struct {
	*core.BaseAggregateRoot

	id        common.UID
	name      string
	createdAt time.Time
	updatedAt time.Time
}

// NewGame - creates a new Game instance with the provided name.
func NewGame(
	name string,
) (Game, error) {
	game := Game{
		BaseAggregateRoot: &core.BaseAggregateRoot{},
		id:                common.NewUID(),
		name:              name,
	}

	game.BaseAggregateRoot.AddEvent(event.GameCreatedEvent{ID: game.ID().String(), Name: name})

	return game, nil
}

func Hydrate(
	id common.UID,
	name string,
	createdAt time.Time,
	updatedAt time.Time,
) Game {
	game := Game{
		BaseAggregateRoot: &core.BaseAggregateRoot{},
		id:                id,
		name:              name,
		createdAt:         createdAt,
		updatedAt:         updatedAt,
	}

	return game
}

func (u *Game) ID() common.UID {
	return u.id
}

func (u *Game) IsEmpty() bool {
	return *u == Game{}
}

// Name returns the name of the game.
func (u *Game) Name() string {
	return u.name
}

func (u *Game) CreatedAt() time.Time {
	return u.createdAt
}

func (u *Game) UpdatedAt() time.Time {
	return u.updatedAt
}

// String returns the string representation of the game.
func (u *Game) String() string {
	return fmt.Sprintf(
		"Game{ID: %s, Name: %s}",
		u.ID(),
		u.Name(),
	)
}

func (u *Game) Events() []mediator.Event {
	return u.BaseAggregateRoot.Events()
}
