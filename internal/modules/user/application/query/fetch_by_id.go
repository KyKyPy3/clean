package query

import (
	"context"
	"fmt"

	"github.com/KyKyPy3/clean/internal/application/core"
	"github.com/KyKyPy3/clean/internal/domain/common"
	"github.com/KyKyPy3/clean/internal/modules/user/application/ports"
	"github.com/KyKyPy3/clean/pkg/logger"
)

const FetchUserByIDKind = "FetchUserByID"

type FetchUserByIDQuery struct {
	ID string
}

func (f FetchUserByIDQuery) Type() core.QueryType {
	return FetchUserByIDKind
}

var _ core.Query = (*FetchUserByIDQuery)(nil)

type FetchUserByID struct {
	storage ports.UserPgStorage
	logger  logger.Logger
}

func NewFetchUserByID(
	storage ports.UserPgStorage,
	logger logger.Logger,
) FetchUserByID {
	return FetchUserByID{
		storage: storage,
		logger:  logger,
	}
}

func (f FetchUserByID) Handle(ctx context.Context, query core.Query) (any, error) {
	fetchByIDQuery, ok := query.(FetchUserByIDQuery)
	if !ok {
		return nil, fmt.Errorf("query type %s: %w", query.Type(), core.ErrUnexpectedQuery)
	}

	id, err := common.ParseUID(fetchByIDQuery.ID)
	if err != nil {
		return nil, err
	}

	user, err := f.storage.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return user, nil
}
