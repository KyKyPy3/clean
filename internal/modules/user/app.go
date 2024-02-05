package user

import (
	"github.com/avito-tech/go-transaction-manager/trm/v2/manager"
	"github.com/labstack/echo/v4"

	"github.com/KyKyPy3/clean/internal/application/core"
	"github.com/KyKyPy3/clean/internal/modules/user/application/command"
	"github.com/KyKyPy3/clean/internal/modules/user/application/ports"
	"github.com/KyKyPy3/clean/internal/modules/user/application/query"
	"github.com/KyKyPy3/clean/internal/modules/user/infrastructure/controller/http/v1"
	"github.com/KyKyPy3/clean/pkg/logger"
	"github.com/KyKyPy3/clean/pkg/mediator"
)

func InitHandlers(
	userPgStorage ports.UserPgStorage,
	mountPoint *echo.Group,
	pubsub *mediator.Mediator,
	trManager *manager.Manager,
	logger logger.Logger,
) {
	userCmdBus := core.NewCommandBus()
	userCmdBus.Register(
		command.CreateUserKind,
		command.NewCreateUser(userPgStorage, trManager, pubsub, logger),
	)
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

	v1.NewUserHandlers(mountPoint, userCmdBus, userQueryBus, logger)
}
