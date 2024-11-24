package query

import (
	"context"
	"fmt"

	"github.com/KyKyPy3/clean/internal/application/core"
	"github.com/KyKyPy3/clean/internal/modules/game/application/ports"
	"github.com/KyKyPy3/clean/internal/modules/game/infrastructure/controller/http/dto"
	"github.com/KyKyPy3/clean/pkg/logger"
)

const FetchGamesKind = "FetchGames"

type FetchGamesQuery struct {
	Limit  int64
	Offset int64
}

func (f FetchGamesQuery) Type() core.QueryType {
	return FetchGamesKind
}

var _ core.Query = (*FetchGamesQuery)(nil)

type FetchGames struct {
	gameStorage ports.GamePgStorage
	userStorage ports.UserViewStorage
	logger      logger.Logger
}

func NewFetchGames(
	gameStorage ports.GamePgStorage,
	userStorage ports.UserViewStorage,
	logger logger.Logger,
) FetchGames {
	return FetchGames{
		gameStorage: gameStorage,
		userStorage: userStorage,
		logger:      logger,
	}
}

func (f FetchGames) Handle(ctx context.Context, query core.Query) (any, error) {
	fetchQuery, ok := query.(FetchGamesQuery)
	if !ok {
		return nil, fmt.Errorf("query type %s: %w", query.Type(), core.ErrUnexpectedQuery)
	}

	games, err := f.gameStorage.Fetch(ctx, fetchQuery.Limit, fetchQuery.Offset)
	if err != nil {
		return nil, err
	}

	gamesDto := make([]dto.GameShortDTO, 0, len(games))
	for _, game := range games {
		user, err := f.userStorage.GetByID(ctx, game.OwnerID())
		if err != nil {
			return nil, err
		}

		gamesDto = append(gamesDto, dto.GameShortDTO{
			ID:        game.ID().String(),
			Name:      game.Name(),
			User:      user.FullName().String(),
			CreatedAt: game.CreatedAt().String(),
			UpdatedAt: game.UpdatedAt().String(),
		})
	}

	return gamesDto, nil
}
