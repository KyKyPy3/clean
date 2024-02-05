package application

import (
	"context"
	"errors"

	"github.com/KyKyPy3/clean/internal/domain/common"
	"github.com/KyKyPy3/clean/internal/domain/core"
	"github.com/KyKyPy3/clean/internal/modules/user/application/ports"
	"github.com/KyKyPy3/clean/pkg/logger"
)

type UniquenessPolicy struct {
	ctx       context.Context
	pgStorage ports.UserViewStorage
	logger    logger.Logger
}

func NewUniquenessPolicy(ctx context.Context, pgStorage ports.UserViewStorage, logger logger.Logger) UniquenessPolicy {
	return UniquenessPolicy{
		ctx:       ctx,
		pgStorage: pgStorage,
		logger:    logger,
	}
}

func (p UniquenessPolicy) IsUnique(email common.Email) (bool, error) {
	_, err := p.pgStorage.GetByEmail(p.ctx, email)
	if err != nil {
		if errors.Is(err, core.ErrNotFound) {
			return true, nil
		}
		return false, err
	}
	return false, nil
}
