package application

import (
	"context"
	"errors"

	"github.com/KyKyPy3/clean/internal/domain/common"
	"github.com/KyKyPy3/clean/internal/domain/core"
	"github.com/KyKyPy3/clean/internal/modules/user/domain/entity"
	"github.com/KyKyPy3/clean/pkg/logger"
)

type UserViewPgStorage interface {
	GetByEmail(ctx context.Context, email common.Email) (entity.User, error)
}

type UniquenessPolicy struct {
	pgStorage UserViewPgStorage
	logger    logger.Logger
}

func NewUniquenessPolicy(pgStorage UserViewPgStorage, logger logger.Logger) UniquenessPolicy {
	return UniquenessPolicy{
		pgStorage: pgStorage,
		logger:    logger,
	}
}

func (p UniquenessPolicy) IsUnique(ctx context.Context, email common.Email) (bool, error) {
	_, err := p.pgStorage.GetByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, core.ErrNotFound) {
			return true, nil
		}
		return false, err
	}
	return false, nil
}
