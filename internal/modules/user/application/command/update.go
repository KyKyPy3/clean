package command

import (
	"context"
	"errors"
	"fmt"

	"github.com/KyKyPy3/clean/internal/application/core"
	"github.com/KyKyPy3/clean/internal/domain/common"
	domain_core "github.com/KyKyPy3/clean/internal/domain/core"
	"github.com/KyKyPy3/clean/internal/modules/user/application/ports"
	"github.com/KyKyPy3/clean/internal/modules/user/domain/entity"
	"github.com/KyKyPy3/clean/internal/modules/user/domain/vo"
	"github.com/KyKyPy3/clean/pkg/logger"
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
	policy   ports.UniquenessPolicer
	manager  ports.TrManager
	mediator ports.Mediator
	logger   logger.Logger
}

func NewUpdateUser(
	storage ports.UserPgStorage,
	manager ports.TrManager,
	mediator ports.Mediator,
	policy ports.UniquenessPolicer,
	logger logger.Logger,
) UpdateUser {
	return UpdateUser{
		storage:  storage,
		manager:  manager,
		mediator: mediator,
		policy:   policy,
		logger:   logger,
	}
}

func (c UpdateUser) Handle(ctx context.Context, command core.Command) (any, error) {
	updateCommand, ok := command.(UpdateUserCommand)
	if !ok {
		return nil, fmt.Errorf("command type %s: %w", command.Type(), core.ErrUnexpectedCommand)
	}

	id, err := common.ParseUID(updateCommand.ID)
	if err != nil {
		return nil, err
	}

	fullname, err := vo.NewFullName(updateCommand.Name, updateCommand.Surname, updateCommand.Middlename)
	if err != nil {
		return nil, err
	}

	email, err := common.NewEmail(updateCommand.Email)
	if err != nil {
		return nil, err
	}

	// TOFIX: add check for email uniq

	err = c.manager.Do(ctx, func(ctx context.Context) error {
		var user entity.User
		user, err = c.storage.GetByID(ctx, id)
		if err != nil && !errors.Is(err, domain_core.ErrNotFound) {
			return err
		}

		if user.IsEmpty() {
			return domain_core.ErrNotFound
		}

		err = user.UpdateEmail(email, c.policy)
		if err != nil {
			return err
		}
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
		return nil, err
	}

	var res interface{}
	return res, nil
}

var _ core.CommandHandler = (*UpdateUser)(nil)
