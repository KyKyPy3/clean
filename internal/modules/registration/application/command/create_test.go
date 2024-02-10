package command_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/KyKyPy3/clean/internal/application/core"
	"github.com/KyKyPy3/clean/internal/domain/common"
	domain_core "github.com/KyKyPy3/clean/internal/domain/core"
	"github.com/KyKyPy3/clean/internal/modules/registration/application/command"
	mocks "github.com/KyKyPy3/clean/mocks/internal_/application/core"
	ports "github.com/KyKyPy3/clean/mocks/internal_/modules/registration/application/ports"
	"github.com/KyKyPy3/clean/pkg/logger"
)

func TestHandleUnsupportedCreateRegistrationCommandError(t *testing.T) {
	log := logger.NewLogger(logger.Config{
		Mode:     "development",
		Level:    "debug",
		Encoding: "json",
	})

	var unsupportedCommandType core.CommandType = "unsupported_crete.verification.command"

	unsupportedCommand := mocks.NewCommand(t)
	unsupportedCommand.On("Type").Return(unsupportedCommandType)

	registrationStorageMock := ports.NewRegistrationPgStorage(t)
	policyMock := ports.NewUniquenessPolicer(t)
	managerMock := ports.NewTrManager(t)
	mediatorMock := ports.NewMediator(t)

	createRegistrationCommandHandler := command.NewCreateRegistration(
		registrationStorageMock,
		policyMock,
		managerMock,
		mediatorMock,
		log,
	)
	_, err := createRegistrationCommandHandler.Handle(context.Background(), unsupportedCommand)

	registrationStorageMock.AssertExpectations(t)
	assert.ErrorIs(t, err, core.ErrUnexpectedCommand)
}

func TestHandleUnsupportedCreateRegistrationEmailError(t *testing.T) {
	log := logger.NewLogger(logger.Config{
		Mode:     "development",
		Level:    "debug",
		Encoding: "json",
	})

	email := "test"
	password := "12345"

	createRegistrationCommand := command.NewCreateRegistrationCommand(email, password)

	registrationStorageMock := ports.NewRegistrationPgStorage(t)
	policyMock := ports.NewUniquenessPolicer(t)
	managerMock := ports.NewTrManager(t)
	mediatorMock := ports.NewMediator(t)

	createRegistrationCommandHandler := command.NewCreateRegistration(
		registrationStorageMock,
		policyMock,
		managerMock,
		mediatorMock,
		log,
	)
	_, err := createRegistrationCommandHandler.Handle(context.Background(), createRegistrationCommand)

	registrationStorageMock.AssertExpectations(t)
	assert.ErrorIs(t, err, common.ErrBadFormat)
}

func TestHandleUnsupportedCreateRegistrationNotUniqError(t *testing.T) {
	log := logger.NewLogger(logger.Config{
		Mode:     "development",
		Level:    "debug",
		Encoding: "json",
	})

	email := "test@mail.com"
	password := "12345"

	createRegistrationCommand := command.NewCreateRegistrationCommand(email, password)

	registrationStorageMock := ports.NewRegistrationPgStorage(t)
	policyMock := ports.NewUniquenessPolicer(t)
	managerMock := ports.NewTrManager(t)
	mediatorMock := ports.NewMediator(t)

	policyMock.On("IsUnique", mock.Anything, mock.Anything).Return(false, nil)

	createRegistrationCommandHandler := command.NewCreateRegistration(
		registrationStorageMock,
		policyMock,
		managerMock,
		mediatorMock,
		log,
	)
	_, err := createRegistrationCommandHandler.Handle(context.Background(), createRegistrationCommand)

	registrationStorageMock.AssertExpectations(t)
	assert.ErrorIs(t, err, domain_core.ErrAlreadyExist)
}

func TestHandleCreateRegistrationCommandSuccess(t *testing.T) {
	log := logger.NewLogger(logger.Config{
		Mode:     "development",
		Level:    "debug",
		Encoding: "json",
	})

	email := "test@gmail.com"
	password := "12345"

	createRegistrationCommand := command.NewCreateRegistrationCommand(email, password)

	registrationStorageMock := ports.NewRegistrationPgStorage(t)
	policyMock := ports.NewUniquenessPolicer(t)
	managerMock := ports.NewTrManager(t)
	mediatorMock := ports.NewMediator(t)

	policyMock.On("IsUnique", mock.Anything, mock.Anything).Return(true, nil)
	managerMock.
		EXPECT().Do(mock.Anything, mock.Anything).
		RunAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
			return fn(ctx)
		})
	mediatorMock.On("Publish", mock.Anything, mock.Anything).Return(nil)
	registrationStorageMock.On("Create", mock.Anything, mock.Anything).Return(nil)

	createRegistrationCommandHandler := command.NewCreateRegistration(
		registrationStorageMock,
		policyMock,
		managerMock,
		mediatorMock,
		log,
	)
	_, err := createRegistrationCommandHandler.Handle(context.Background(), createRegistrationCommand)

	registrationStorageMock.AssertExpectations(t)
	assert.NoError(t, err)
}
