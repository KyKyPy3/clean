package query

import (
	"context"
	"fmt"

	"github.com/KyKyPy3/clean/internal/application/core"
	"github.com/KyKyPy3/clean/internal/modules/user/application/ports"
	"github.com/KyKyPy3/clean/pkg/logger"
)

const FetchUsersKind = "FetchUsers"

type FetchUsersQuery struct {
	Limit  int64
	Offset int64
}

func (f FetchUsersQuery) Type() core.QueryType {
	return FetchUsersKind
}

var _ core.Query = (*FetchUsersQuery)(nil)

type FetchUsers struct {
	storage ports.UserPgStorage
	logger  logger.Logger
}

func NewFetchUsers(
	storage ports.UserPgStorage,
	logger logger.Logger,
) FetchUsers {
	return FetchUsers{
		storage: storage,
		logger:  logger,
	}
}

func (f FetchUsers) Handle(ctx context.Context, query core.Query) (any, error) {
	fetchQuery, ok := query.(FetchUsersQuery)
	if !ok {
		return nil, fmt.Errorf("query type %s: %w", query.Type(), core.ErrUnexpectedQuery)
	}

	users, err := f.storage.Fetch(ctx, fetchQuery.Limit, fetchQuery.Offset)
	if err != nil {
		return nil, err
	}

	return users, nil
}
