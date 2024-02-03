package command

import (
	"context"
	"fmt"

	"github.com/KyKyPy3/clean/internal/application/core"
	"github.com/KyKyPy3/clean/internal/domain/common"
	"github.com/KyKyPy3/clean/internal/modules/registration/application/ports"
	"github.com/KyKyPy3/clean/internal/modules/registration/domain/entity"
	"github.com/KyKyPy3/clean/pkg/logger"
)

const CreateRegistrationKind = "CreateRegistration"

type CreateRegistrationCommand struct {
	Email    string
	Password string
}

func NewCreateRegistrationCommand(email string, password string) CreateRegistrationCommand {
	return CreateRegistrationCommand{
		Email:    email,
		Password: password,
	}
}

func (c CreateRegistrationCommand) Type() core.CommandType {
	return CreateRegistrationKind
}

var _ core.Command = (*CreateRegistrationCommand)(nil)

type CreateRegistration struct {
	storage  ports.RegistrationPgStorage
	policy   ports.UniquenessPolicer
	manager  ports.TrManager
	mediator ports.Mediator
	logger   logger.Logger
}

func NewCreateRegistration(
	storage ports.RegistrationPgStorage,
	policy ports.UniquenessPolicer,
	manager ports.TrManager,
	mediator ports.Mediator,
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

func (c CreateRegistration) Handle(ctx context.Context, command core.Command) error {
	createCommand, ok := command.(CreateRegistrationCommand)
	if !ok {
		return fmt.Errorf("command type %s: %w", command.Type(), core.ErrUnexpectedCommand)
	}

	email, err := common.NewEmail(createCommand.Email)
	if err != nil {
		return err
	}

	reg, err := entity.NewRegistration(ctx, email, createCommand.Password, c.policy)
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

var _ core.CommandHandler = (*CreateRegistration)(nil)
