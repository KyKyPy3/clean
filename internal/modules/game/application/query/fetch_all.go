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

	games, fetchErr := f.gameStorage.Fetch(ctx, fetchQuery.Limit, fetchQuery.Offset)
	if fetchErr != nil {
		return nil, fmt.Errorf("fetch games error: %w", fetchErr)
	}

	gamesDto := make([]dto.GameShortDTO, 0, len(games))
	for _, game := range games {
		user, userErr := f.userStorage.GetByID(ctx, game.OwnerID())
		if userErr != nil {
			return nil, fmt.Errorf("get user error for game %s: %w", game.ID(), userErr)
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
