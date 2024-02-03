package command

import (
	"context"
	"errors"
	"fmt"

	"github.com/avito-tech/go-transaction-manager/trm/v2/manager"

	"github.com/KyKyPy3/clean/internal/application/core"
	"github.com/KyKyPy3/clean/internal/domain/common"
	domain_core "github.com/KyKyPy3/clean/internal/domain/core"
	"github.com/KyKyPy3/clean/internal/modules/user/application/ports"
	"github.com/KyKyPy3/clean/internal/modules/user/domain/value_object"
	"github.com/KyKyPy3/clean/pkg/logger"
	"github.com/KyKyPy3/clean/pkg/mediator"
)

const UpdateUserKind = "UpdateUser"

type UpdateUserCommand struct {
	ID         string
	Name       string
	Surname    string
	Middlename string
	Email      string
}

func (c UpdateUserCommand) Type() core.CommandType {
	return UpdateUserKind
}

var _ core.Command = (*UpdateUserCommand)(nil)

type UpdateUser struct {
	storage  ports.UserPgStorage
	manager  *manager.Manager
	mediator *mediator.Mediator
	logger   logger.Logger
}

func NewUpdateUser(
	storage ports.UserPgStorage,
	manager *manager.Manager,
	mediator *mediator.Mediator,
	logger logger.Logger,
) UpdateUser {
	return UpdateUser{
		storage:  storage,
		manager:  manager,
		mediator: mediator,
		logger:   logger,
	}
}

func (c UpdateUser) Handle(ctx context.Context, command core.Command) error {
	updateCommand, ok := command.(UpdateUserCommand)
	if !ok {
		return fmt.Errorf("command type %s: %w", command.Type(), core.ErrUnexpectedCommand)
	}

	id, err := common.ParseUID(updateCommand.ID)
	if err != nil {
		return err
	}

	fullname, err := value_object.NewFullName(updateCommand.Name, updateCommand.Surname, updateCommand.Middlename)
	if err != nil {
		return err
	}

	email, err := common.NewEmail(updateCommand.Email)
	if err != nil {
		return err
	}

	err = c.manager.Do(ctx, func(ctx context.Context) error {
		user, err := c.storage.GetByID(ctx, id)
		if err != nil && !errors.Is(err, domain_core.ErrNotFound) {
			return err
		}

		if user.IsEmpty() {
			return domain_core.ErrNotFound
		}

		user.UpdateEmail(email)
		user.UpdateFullName(fullname)

		err = c.storage.Update(ctx, user)
		if err != nil {
			return err
		}

		err = c.mediator.Publish(ctx, user.Events()...)
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

var _ core.CommandHandler = (*UpdateUser)(nil)
