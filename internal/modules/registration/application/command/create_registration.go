package command

import (
	"context"

	"github.com/avito-tech/go-transaction-manager/trm/v2/manager"

	"github.com/KyKyPy3/clean/internal/application/core"
	"github.com/KyKyPy3/clean/internal/domain/common"
	"github.com/KyKyPy3/clean/internal/modules/registration/application/ports"
	"github.com/KyKyPy3/clean/internal/modules/registration/domain/entity"
	"github.com/KyKyPy3/clean/pkg/logger"
	"github.com/KyKyPy3/clean/pkg/mediator"
)

type CreateRegistrationCommand struct {
	Email string
}

type CreateRegistrationHandler core.CommandHandler[CreateRegistrationCommand]

type UniquenessPolicy interface {
	IsUnique(ctx context.Context, email common.Email) (bool, error)
}

type CreateRegistration struct {
	storage  ports.RegistrationPgStorage
	policy   UniquenessPolicy
	manager  *manager.Manager
	mediator *mediator.Mediator
	logger   logger.Logger
}

func NewCreateRegistration(
	storage ports.RegistrationPgStorage,
	policy UniquenessPolicy,
	manager *manager.Manager,
	mediator *mediator.Mediator,
	logger logger.Logger,
) CreateRegistration {
	return CreateRegistration{
		storage:  storage,
		policy:   policy,
		manager:  manager,
		mediator: mediator,
		logger:   logger,
	}
}

func (c CreateRegistration) Handle(ctx context.Context, command CreateRegistrationCommand) error {
	email, err := common.NewEmail(command.Email)
	if err != nil {
		return err
	}

	reg, err := entity.NewRegistration(ctx, email, c.policy)
	if err != nil {
		return err
	}

	err = c.manager.Do(ctx, func(ctx context.Context) error {
		err := c.storage.Create(ctx, reg)
		if err != nil {
			return err
		}

		err = c.mediator.Publish(ctx, reg.Events()...)
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
