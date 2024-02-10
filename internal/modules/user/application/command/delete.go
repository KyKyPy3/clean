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
	"github.com/KyKyPy3/clean/pkg/logger"
)

const DeleteUserKind = "DeleteUser"

type DeleteUserCommand struct {
	ID string
}

func (c DeleteUserCommand) Type() core.CommandType {
	return DeleteUserKind
}

var _ core.Command = (*DeleteUserCommand)(nil)

type DeleteUser struct {
	storage  ports.UserPgStorage
	manager  ports.TrManager
	mediator ports.Mediator
	logger   logger.Logger
}

func NewDeleteUser(
	storage ports.UserPgStorage,
	manager ports.TrManager,
	mediator ports.Mediator,
	logger logger.Logger,
) DeleteUser {
	return DeleteUser{
		storage:  storage,
		manager:  manager,
		mediator: mediator,
		logger:   logger,
	}
}

func (c DeleteUser) Handle(ctx context.Context, command core.Command) (any, error) {
	deleteCommand, ok := command.(DeleteUserCommand)
	if !ok {
		return nil, fmt.Errorf("command type %s: %w", command.Type(), core.ErrUnexpectedCommand)
	}

	id, err := common.ParseUID(deleteCommand.ID)
	if err != nil {
		return nil, err
	}

	err = c.manager.Do(ctx, func(ctx context.Context) error {
		var user entity.User
		user, err = c.storage.GetByID(ctx, id)
		if err != nil && !errors.Is(err, domain_core.ErrNotFound) {
			return err
		}

		if user.IsEmpty() {
			return domain_core.ErrNotFound
		}

		err = c.storage.Delete(ctx, id)
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

var _ core.CommandHandler = (*DeleteUser)(nil)
