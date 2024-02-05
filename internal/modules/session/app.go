package session

import (
	"github.com/labstack/echo/v4"

	"github.com/KyKyPy3/clean/internal/application/core"
	"github.com/KyKyPy3/clean/internal/infrastructure/config"
	"github.com/KyKyPy3/clean/internal/modules/session/application/command"
	"github.com/KyKyPy3/clean/internal/modules/session/application/ports"
	handlers "github.com/KyKyPy3/clean/internal/modules/session/infrastructure/controller/http/v1"
	"github.com/KyKyPy3/clean/pkg/jwt"
	"github.com/KyKyPy3/clean/pkg/logger"
)

func InitHandlers(
	userPgStorage ports.UserPgStorage,
	sessionRedisStorage ports.SessionRedisStorage,
	publicMountPoint *echo.Group,
	privateMountPoint *echo.Group,
	cfg *config.Config,
	jwt *jwt.JWT,
	logger logger.Logger,
) {
	regCmdBus := core.NewCommandBus()
	regCmdBus.Register(
		command.LoginUserKind,
		command.NewLoginUser(userPgStorage, sessionRedisStorage, logger),
	)
	regCmdBus.Register(
		command.LogoutUserKind,
		command.NewLogoutUser(sessionRedisStorage, logger),
	)

	userQueryBus := core.NewQueryBus()

	handlers.NewAuthHandlers(publicMountPoint, privateMountPoint, regCmdBus, userQueryBus, cfg, jwt, logger)
}
