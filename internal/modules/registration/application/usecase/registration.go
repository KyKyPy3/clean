package usecase

import (
	"context"
	"github.com/KyKyPy3/clean/internal/modules/registration/domain/event"
	"github.com/KyKyPy3/clean/pkg/outbox"
	"time"

	"github.com/avito-tech/go-transaction-manager/trm/v2/manager"

	"github.com/KyKyPy3/clean/internal/modules/registration/domain/entity"
	v1 "github.com/KyKyPy3/clean/internal/modules/registration/infrastructure/controller/http/v1"
	"github.com/KyKyPy3/clean/pkg/logger"
)

const (
	requestTimeout = time.Second * 5
)

type RegistrationService interface {
	Create(ctx context.Context, data entity.Registration) error
}

type registrationUsecase struct {
	registrationService RegistrationService
	logger              logger.Logger
	manager             *manager.Manager
	outbox              outbox.Manager
}

func NewRegistrationUsecase(registrationService RegistrationService, manager *manager.Manager, outbox outbox.Manager, logger logger.Logger) v1.RegistrationUsecase {
	return &registrationUsecase{
		registrationService: registrationService,
		manager:             manager,
		logger:              logger,
		outbox:              outbox,
	}
}

func (r *registrationUsecase) Create(ctx context.Context, data entity.Registration) error {
	ctx, cancel := context.WithTimeout(ctx, requestTimeout)
	defer cancel()

	err := r.manager.Do(ctx, func(ctx context.Context) error {
		err := r.registrationService.Create(ctx, data)
		if err != nil {
			return err
		}

		err = r.outbox.Publish(ctx, event.RegistrationCreatedEvent{ID: "1", Email: data.GetEmail()})
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}
