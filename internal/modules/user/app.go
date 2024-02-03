package user

import (
	trmsqlx "github.com/avito-tech/go-transaction-manager/drivers/sqlx/v2"
	"github.com/avito-tech/go-transaction-manager/trm/v2/manager"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"

	"github.com/KyKyPy3/clean/internal/application/core"
	"github.com/KyKyPy3/clean/internal/modules/user/application/command"
	"github.com/KyKyPy3/clean/internal/modules/user/application/query"
	"github.com/KyKyPy3/clean/internal/modules/user/infrastructure/controller/http/v1"
	"github.com/KyKyPy3/clean/internal/modules/user/infrastructure/gateway/postgres"
	"github.com/KyKyPy3/clean/pkg/logger"
	"github.com/KyKyPy3/clean/pkg/mediator"
)

func InitUserHandlers(
	pgClient *sqlx.DB,
	mountPoint *echo.Group,
	pubsub *mediator.Mediator,
	trManager *manager.Manager,
	logger logger.Logger,
) {
	userPgStorage := postgres.NewUserPgStorage(pgClient, trmsqlx.DefaultCtxGetter, logger)
	//userRedisStorage := userRedis.NewUserRedisStorage(a.redisClient, logger)

	userCmdBus := core.NewCommandBus()
	userCmdBus.Register(
		command.CreateUserKind,
		command.NewCreateUser(userPgStorage, trManager, pubsub, logger),
	)
	userQueryBus := core.NewQueryBus()
	userQueryBus.Register(
		query.FetchUsersKind,
		query.NewFetchUsers(userPgStorage, logger),
	)

	v1.NewUserHandlers(mountPoint, userCmdBus, userQueryBus, logger)
}
