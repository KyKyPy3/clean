package postgres

import (
	"time"

	"github.com/KyKyPy3/clean/internal/domain/common"
	"github.com/KyKyPy3/clean/internal/modules/game/domain/entity"
)

// DBGame Database game representation.
type DBGame struct {
	ID        string    `db:"id"`
	Name      string    `db:"name"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

// GameFromDB Convert database game model to domain model.
func GameFromDB(dbGame DBGame) (entity.Game, error) {
	entityID, err := common.ParseUID(dbGame.ID)
	if err != nil {
		return entity.Game{}, err
	}

	game := entity.Hydrate(entityID, dbGame.Name, dbGame.CreatedAt, dbGame.UpdatedAt)

	return game, nil
}

// GameToDB Convert domain game model to database model.
func GameToDB(game entity.Game) DBGame {
	return DBGame{
		ID:   game.ID().String(),
		Name: game.Name(),
	}
}
