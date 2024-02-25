package query

import (
	"context"
	"fmt"

	"github.com/KyKyPy3/clean/internal/application/core"
	"github.com/KyKyPy3/clean/internal/modules/game/application/ports"
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
	storage ports.GamePgStorage
	logger  logger.Logger
}

func NewFetchGames(
	storage ports.GamePgStorage,
	logger logger.Logger,
) FetchGames {
	return FetchGames{
		storage: storage,
		logger:  logger,
	}
}

func (f FetchGames) Handle(ctx context.Context, query core.Query) (any, error) {
	fetchQuery, ok := query.(FetchGamesQuery)
	if !ok {
		return nil, fmt.Errorf("query type %s: %w", query.Type(), core.ErrUnexpectedQuery)
	}

	games, err := f.storage.Fetch(ctx, fetchQuery.Limit, fetchQuery.Offset)
	if err != nil {
		return nil, err
	}

	return games, nil
}
