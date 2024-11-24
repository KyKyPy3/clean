package entity

import (
	"fmt"
	"reflect"
	"time"

	"github.com/KyKyPy3/clean/internal/domain/common"
	"github.com/KyKyPy3/clean/internal/domain/core"
	"github.com/KyKyPy3/clean/internal/modules/game/domain/event"
	"github.com/KyKyPy3/clean/pkg/mediator"
)

// Game struct.
type Game struct {
	*core.BaseAggregateRoot

	id          common.UID
	name        string
	ownerID     common.UID
	gameMembers []GameMember
	createdAt   time.Time
	updatedAt   time.Time
}

// NewGame - creates a new Game instance with the provided name.
func NewGame(
	name string,
	ownerID common.UID,
) (Game, error) {
	game := Game{
		BaseAggregateRoot: &core.BaseAggregateRoot{},
		id:                common.NewUID(),
		name:              name,
		ownerID:           ownerID,
	}

	game.BaseAggregateRoot.AddEvent(event.GameCreatedEvent{ID: game.ID().String(), Name: name})

	return game, nil
}

func Hydrate(
	id common.UID,
	name string,
	ownerID common.UID,
	gameMembers []GameMember,
	createdAt time.Time,
	updatedAt time.Time,
) Game {
	game := Game{
		BaseAggregateRoot: &core.BaseAggregateRoot{},
		id:                id,
		name:              name,
		ownerID:           ownerID,
		gameMembers:       gameMembers,
		createdAt:         createdAt,
		updatedAt:         updatedAt,
	}

	return game
}

func (g *Game) ID() common.UID {
	return g.id
}

func (g *Game) IsEmpty() bool {
	return reflect.DeepEqual(*g, Game{})
}

// Name returns the name of the game.
func (g *Game) Name() string {
	return g.name
}

func (g *Game) OwnerID() common.UID {
	return g.ownerID
}

func (g *Game) GameMembers() []GameMember {
	return g.gameMembers
}

func (g *Game) AddMember(memberID common.UID, name string) error {
	newMember, err := NewGroupMember(name, memberID, g.id)
	if err != nil {
		return err
	}

	g.gameMembers = append(g.gameMembers, newMember)

	return nil
}

func (g *Game) CreatedAt() time.Time {
	return g.createdAt
}

func (g *Game) UpdatedAt() time.Time {
	return g.updatedAt
}

// String returns the string representation of the game.
// TODO: add additional fields
func (g *Game) String() string {
	return fmt.Sprintf(
		"Game{ID: %s, Name: %s}",
		g.ID(),
		g.Name(),
	)
}

func (g *Game) Events() []mediator.Event {
	return g.BaseAggregateRoot.Events()
}
