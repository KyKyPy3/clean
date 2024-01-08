package service

import (
	"context"
	"errors"
	"github.com/KyKyPy3/clean/internal/domain"
	"github.com/KyKyPy3/clean/internal/modules/registration/domain/entity"
	"github.com/KyKyPy3/clean/pkg/logger"

	"github.com/KyKyPy3/clean/internal/domain/common"
)

type RegistrationViewPgStorage interface {
	GetByEmail(ctx context.Context, email common.Email) (entity.Registration, error)
}

type UniquenessPolicy struct {
	pgStorage RegistrationViewPgStorage
	logger    logger.Logger
}

func NewUniquenessPolicy(pgStorage RegistrationViewPgStorage, logger logger.Logger) UniquenessPolicy {
	return UniquenessPolicy{
		pgStorage: pgStorage,
		logger:    logger,
	}
}

func (p UniquenessPolicy) IsUnique(ctx context.Context, email common.Email) (bool, error) {
	_, err := p.pgStorage.GetByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return true, nil
		}
		return false, err
	}
	return false, nil
}
