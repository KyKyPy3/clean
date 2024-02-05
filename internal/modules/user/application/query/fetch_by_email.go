package query

import (
	"context"
	"fmt"

	"github.com/KyKyPy3/clean/internal/application/core"
	"github.com/KyKyPy3/clean/internal/domain/common"
	"github.com/KyKyPy3/clean/internal/modules/user/application/ports"
	"github.com/KyKyPy3/clean/pkg/logger"
)

const FetchUserByEmailKind = "FetchUserByEmail"

type FetchUserByEmailQuery struct {
	Email string
}

func (f FetchUserByEmailQuery) Type() core.QueryType {
	return FetchUserByEmailKind
}

var _ core.Query = (*FetchUserByEmailQuery)(nil)

type FetchUserByEmail struct {
	storage ports.UserPgStorage
	logger  logger.Logger
}

func NewFetchUserByEmail(
	storage ports.UserPgStorage,
	logger logger.Logger,
) FetchUserByEmail {
	return FetchUserByEmail{
		storage: storage,
		logger:  logger,
	}
}

func (f FetchUserByEmail) Handle(ctx context.Context, query core.Query) (any, error) {
	fetchByEmailQuery, ok := query.(FetchUserByEmailQuery)
	if !ok {
		return nil, fmt.Errorf("query type %s: %w", query.Type(), core.ErrUnexpectedQuery)
	}

	email, err := common.NewEmail(fetchByEmailQuery.Email)
	if err != nil {
		return nil, err
	}

	user, err := f.storage.GetByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	return &user, nil
}
