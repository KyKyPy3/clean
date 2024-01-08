package service

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"

	"github.com/KyKyPy3/clean/internal/domain"
	"github.com/KyKyPy3/clean/internal/domain/common"
	"github.com/KyKyPy3/clean/internal/modules/registration/application/usecase"
	"github.com/KyKyPy3/clean/internal/modules/registration/domain/entity"
	"github.com/KyKyPy3/clean/pkg/logger"
)

type RegistrationPgStorage interface {
	Create(ctx context.Context, user entity.Registration) error
	GetByEmail(ctx context.Context, email common.Email) (entity.Registration, error)
}

type UniquenessPolicer interface {
	IsUnique(ctx context.Context, email common.Email) (bool, error)
}

type registrationService struct {
	pgStorage RegistrationPgStorage
	policy    UniquenessPolicer
	logger    logger.Logger
	tracer    trace.Tracer
}

func NewRegistrationService(pgStorage RegistrationPgStorage, policy UniquenessPolicer, logger logger.Logger) usecase.RegistrationService {
	return &registrationService{
		pgStorage: pgStorage,
		policy:    policy,
		logger:    logger,
		tracer:    otel.Tracer(""),
	}
}

func (r *registrationService) Create(ctx context.Context, data entity.Registration) error {
	ctx, span := r.tracer.Start(ctx, "userService.Create")
	defer span.End()

	ok, err := r.policy.IsUnique(ctx, data.GetEmail())
	if err != nil {
		return fmt.Errorf("failed to check uniqueness of email on registration, err: %w", err)
	}

	if !ok {
		return domain.ErrAlreadyExist
	}

	err = r.pgStorage.Create(ctx, data)
	if err != nil {
		return err
	}

	return nil
}
