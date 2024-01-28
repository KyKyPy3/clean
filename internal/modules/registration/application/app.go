package application

import (
	"github.com/avito-tech/go-transaction-manager/trm/v2/manager"

	"github.com/KyKyPy3/clean/internal/modules/registration/application/command"
	"github.com/KyKyPy3/clean/internal/modules/registration/application/ports"
	"github.com/KyKyPy3/clean/pkg/logger"
	"github.com/KyKyPy3/clean/pkg/mediator"
)

type Application struct {
	Commands Commands
	Queries  Queries
}

type Commands struct {
	CreateRegistration  command.CreateRegistrationHandler
	ConfirmRegistration command.ConfirmRegistration
}

type Queries struct {
}

func NewApplication(
	storage ports.RegistrationPgStorage,
	policy ports.UniquenessPolicer,
	manager *manager.Manager,
	mediator *mediator.Mediator,
	logger logger.Logger,
) *Application {
	commands := Commands{
		CreateRegistration:  command.NewCreateRegistration(storage, policy, manager, mediator, logger),
		ConfirmRegistration: command.NewConfirmRegistration(storage, mediator, manager, logger),
	}

	return &Application{
		Commands: commands,
	}
}
