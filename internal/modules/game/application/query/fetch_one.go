package query

import (
	"context"
	"fmt"

	"github.com/KyKyPy3/clean/internal/application/core"
	"github.com/KyKyPy3/clean/internal/domain/common"
	"github.com/KyKyPy3/clean/internal/modules/game/application/ports"
	"github.com/KyKyPy3/clean/internal/modules/game/infrastructure/controller/http/dto"
	"github.com/KyKyPy3/clean/pkg/logger"
)

const FetchGameKind = "FetchGame"

type FetchGameQuery struct {
	ID string
}

func (f FetchGameQuery) Type() core.QueryType {
	return FetchGameKind
}

var _ core.Query = (*FetchGameQuery)(nil)

type FetchGame struct {
	gameStorage ports.GamePgStorage
	userStorage ports.UserViewStorage
	logger      logger.Logger
}

func NewFetchGame(
	gameStorage ports.GamePgStorage,
	userStorage ports.UserViewStorage,
	logger logger.Logger,
) FetchGame {
	return FetchGame{
		gameStorage: gameStorage,
		userStorage: userStorage,
		logger:      logger,
	}
}

func (f FetchGame) Handle(ctx context.Context, query core.Query) (any, error) {
	fetchByIDQuery, ok := query.(FetchGameQuery)
	if !ok {
		return nil, fmt.Errorf("query type %s: %w", query.Type(), core.ErrUnexpectedQuery)
	}

	id, err := common.ParseUID(fetchByIDQuery.ID)
	if err != nil {
		return nil, err
	}

	game, err := f.gameStorage.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	user, err := f.userStorage.GetByID(ctx, game.OwnerID())
	if err != nil {
		return nil, err
	}

	return dto.GameDTO{
		ID:        game.ID().String(),
		Name:      game.Name(),
		User:      user.FullName().String(),
		CreatedAt: game.CreatedAt().String(),
		UpdatedAt: game.UpdatedAt().String(),
	}, nil
}
