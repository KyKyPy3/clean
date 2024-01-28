package command

import (
	"context"
	"errors"
	"github.com/KyKyPy3/clean/pkg/mediator"

	"github.com/avito-tech/go-transaction-manager/trm/v2/manager"

	"github.com/KyKyPy3/clean/internal/application/core"
	"github.com/KyKyPy3/clean/internal/domain/common"
	domain_core "github.com/KyKyPy3/clean/internal/domain/core"
	"github.com/KyKyPy3/clean/internal/modules/registration/application/ports"
	"github.com/KyKyPy3/clean/pkg/logger"
)

type ConfirmRegistrationCommand struct {
	ID string
}

type ConfirmRegistrationHandler core.CommandHandler[ConfirmRegistrationCommand]

type ConfirmRegistration struct {
	manager  *manager.Manager
	storage  ports.RegistrationPgStorage
	mediator *mediator.Mediator
	logger   logger.Logger
}

func NewConfirmRegistration(storage ports.RegistrationPgStorage, mediator *mediator.Mediator, manager *manager.Manager, logger logger.Logger) ConfirmRegistration {
	return ConfirmRegistration{
		storage:  storage,
		manager:  manager,
		mediator: mediator,
		logger:   logger,
	}
}

func (c ConfirmRegistration) Handle(ctx context.Context, command ConfirmRegistrationCommand) error {
	id, err := common.ParseUID(command.ID)
	if err != nil {
		return err
	}

	err = c.manager.Do(ctx, func(ctx context.Context) error {
		reg, err := c.storage.GetByID(ctx, id)
		if err != nil {
			return err
		}

		err = reg.Verify()
		if err != nil {
			if errors.Is(err, domain_core.ErrNoChanges) {
				return nil
			}

			return err
		}

		err = c.storage.Update(ctx, reg)
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
