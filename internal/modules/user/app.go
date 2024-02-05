package user

import (
	"context"

	"github.com/avito-tech/go-transaction-manager/trm/v2/manager"
	"github.com/labstack/echo/v4"

	"github.com/KyKyPy3/clean/internal/application/core"
	reg_event "github.com/KyKyPy3/clean/internal/modules/registration/domain/event"
	"github.com/KyKyPy3/clean/internal/modules/user/application"
	"github.com/KyKyPy3/clean/internal/modules/user/application/command"
	"github.com/KyKyPy3/clean/internal/modules/user/application/event"
	"github.com/KyKyPy3/clean/internal/modules/user/application/ports"
	"github.com/KyKyPy3/clean/internal/modules/user/application/query"
	"github.com/KyKyPy3/clean/internal/modules/user/infrastructure/controller/http/v1"
	"github.com/KyKyPy3/clean/pkg/logger"
	"github.com/KyKyPy3/clean/pkg/mediator"
)

func InitHandlers(
	ctx context.Context,
	userPgStorage ports.UserPgStorage,
	mountPoint *echo.Group,
	pubsub *mediator.Mediator,
	trManager *manager.Manager,
	logger logger.Logger,
) {
	regUniqPolicy := application.NewUniquenessPolicy(ctx, userPgStorage, logger)
	userCmdBus := core.NewCommandBus()
	userCmdBus.Register(
		command.DeleteUserKind,
		command.NewDeleteUser(userPgStorage, trManager, pubsub, logger),
	)
	userQueryBus := core.NewQueryBus()
	userQueryBus.Register(
		query.FetchUsersKind,
		query.NewFetchUsers(userPgStorage, logger),
	)
	userQueryBus.Register(
		query.FetchUserByIDKind,
		query.NewFetchUserByID(userPgStorage, logger),
	)

	pubsub.Subscribe(reg_event.RegistrationVerified, event.NewRegistrationVerified(logger, userPgStorage, regUniqPolicy).Handle)

	v1.NewUserHandlers(mountPoint, userCmdBus, userQueryBus, logger)
}
