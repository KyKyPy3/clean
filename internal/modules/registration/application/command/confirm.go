package command

import (
	"context"
	"errors"
	"fmt"

	"github.com/KyKyPy3/clean/internal/application/core"
	"github.com/KyKyPy3/clean/internal/domain/common"
	domain_core "github.com/KyKyPy3/clean/internal/domain/core"
	"github.com/KyKyPy3/clean/internal/modules/registration/application/ports"
	"github.com/KyKyPy3/clean/pkg/logger"
)

const ConfirmRegistrationKind = "ConfirmRegistration"

type ConfirmRegistrationCommand struct {
	ID string
}

func (c ConfirmRegistrationCommand) Type() core.CommandType {
	return ConfirmRegistrationKind
}

var _ core.Command = (*ConfirmRegistrationCommand)(nil)

type ConfirmRegistration struct {
	manager  ports.TrManager
	storage  ports.RegistrationPgStorage
	mediator ports.Mediator
	logger   logger.Logger
}

func NewConfirmRegistration(
	storage ports.RegistrationPgStorage,
	mediator ports.Mediator,
	manager ports.TrManager,
	logger logger.Logger,
) ConfirmRegistration {
	return ConfirmRegistration{
		storage:  storage,
		manager:  manager,
		mediator: mediator,
		logger:   logger,
	}
}

func (c ConfirmRegistration) Handle(ctx context.Context, command core.Command) (any, error) {
	confirmCommand, ok := command.(ConfirmRegistrationCommand)
	if !ok {
		return nil, fmt.Errorf("command type %s: %w", command.Type(), core.ErrUnexpectedCommand)
	}

	id, err := common.ParseUID(confirmCommand.ID)
	if err != nil {
		return nil, err
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
		return nil, err
	}

	return nil, nil
}

var _ core.CommandHandler = (*ConfirmRegistration)(nil)
