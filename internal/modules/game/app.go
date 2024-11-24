package game

import (
	"context"

	"github.com/avito-tech/go-transaction-manager/trm/v2/manager"
	"github.com/labstack/echo/v4"

	"github.com/KyKyPy3/clean/internal/application/core"
	"github.com/KyKyPy3/clean/internal/modules/game/application/command"
	"github.com/KyKyPy3/clean/internal/modules/game/application/ports"
	"github.com/KyKyPy3/clean/internal/modules/game/application/query"
	handlers "github.com/KyKyPy3/clean/internal/modules/game/infrastructure/controller/http/v1"
	"github.com/KyKyPy3/clean/pkg/logger"
	"github.com/KyKyPy3/clean/pkg/mediator"
)

func InitHandlers(
	_ context.Context,
	gameStorage ports.GamePgStorage,
	userStorage ports.UserViewStorage,
	mountPoint *echo.Group,
	pubsub *mediator.Mediator,
	trManager *manager.Manager,
	logger logger.Logger,
) {
	gameCmdBus := core.NewCommandBus()
	gameCmdBus.Register(
		command.CreateGameKind,
		command.NewCreateGame(gameStorage, trManager, pubsub, logger),
	)

	gameQueryBus := core.NewQueryBus()
	gameQueryBus.Register(
		query.FetchGamesKind,
		query.NewFetchGames(gameStorage, userStorage, logger),
	)

	handlers.NewGameHandlers(mountPoint, gameCmdBus, gameQueryBus, logger)
}
