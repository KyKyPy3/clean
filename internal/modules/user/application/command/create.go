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
	"github.com/KyKyPy3/clean/internal/modules/user/domain/entity"
	"github.com/KyKyPy3/clean/internal/modules/user/domain/value_object"
	"github.com/KyKyPy3/clean/pkg/logger"
	"github.com/KyKyPy3/clean/pkg/mediator"
)

const CreateUserKind = "CreateUser"

type CreateUserCommand struct {
	Name       string
	Surname    string
	Middlename string
	Email      string
	Password   string
}

func (c CreateUserCommand) Type() core.CommandType {
	return CreateUserKind
}

var _ core.Command = (*CreateUserCommand)(nil)

type CreateUser struct {
	storage  ports.UserPgStorage
	manager  *manager.Manager
	mediator *mediator.Mediator
	logger   logger.Logger
}

func NewCreateUser(
	storage ports.UserPgStorage,
	manager *manager.Manager,
	mediator *mediator.Mediator,
	logger logger.Logger,
) CreateUser {
	return CreateUser{
		storage:  storage,
		manager:  manager,
		mediator: mediator,
		logger:   logger,
	}
}

func (c CreateUser) Handle(ctx context.Context, command core.Command) error {
	createCommand, ok := command.(CreateUserCommand)
	if !ok {
		return fmt.Errorf("command type %s: %w", command.Type(), core.ErrUnexpectedCommand)
	}

	fullname, err := value_object.NewFullName(createCommand.Name, createCommand.Surname, createCommand.Middlename)
	if err != nil {
		return err
	}

	email, err := common.NewEmail(createCommand.Email)
	if err != nil {
		return err
	}

	user, err := entity.NewUser(fullname, email, createCommand.Password)
	if err != nil {
		return err
	}

	err = c.manager.Do(ctx, func(ctx context.Context) error {
		existed, err := c.storage.GetByEmail(ctx, email)
		if err != nil && !errors.Is(err, domain_core.ErrNotFound) {
			return err
		}

		if !existed.IsEmpty() {
			return domain_core.ErrAlreadyExist
		}

		err = c.storage.Create(ctx, user)
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

var _ core.CommandHandler = (*CreateUser)(nil)
