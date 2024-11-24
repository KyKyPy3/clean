package command

import (
	"context"
	"fmt"

	"github.com/KyKyPy3/clean/internal/application/core"
	"github.com/KyKyPy3/clean/internal/domain/common"
	"github.com/KyKyPy3/clean/internal/modules/game/application/ports"
	"github.com/KyKyPy3/clean/internal/modules/game/domain/entity"
	"github.com/KyKyPy3/clean/pkg/logger"
)

const CreateGameKind = "CreateGame"

type CreateGameCommand struct {
	Name   string
	UserID string
}

func NewCreateGameCommand(name, userID string) CreateGameCommand {
	return CreateGameCommand{
		Name:   name,
		UserID: userID,
	}
}

func (c CreateGameCommand) Type() core.CommandType {
	return CreateGameKind
}

var _ core.Command = (*CreateGameCommand)(nil)

type CreateGame struct {
	storage  ports.GamePgStorage
	manager  ports.TrManager
	mediator ports.Mediator
	logger   logger.Logger
}

func NewCreateGame(
	storage ports.GamePgStorage,
	manager ports.TrManager,
	mediator ports.Mediator,
	logger logger.Logger,
) CreateGame {
	return CreateGame{
		storage:  storage,
		manager:  manager,
		mediator: mediator,
		logger:   logger,
	}
}

func (c CreateGame) Handle(ctx context.Context, command core.Command) (any, error) {
	createCommand, ok := command.(CreateGameCommand)
	if !ok {
		return nil, fmt.Errorf("command type %s: %w", command.Type(), core.ErrUnexpectedCommand)
	}

	userID, err := common.ParseUID(createCommand.UserID)
	if err != nil {
		return nil, err
	}

	game, err := entity.NewGame(createCommand.Name, userID)
	if err != nil {
		return nil, err
	}

	err = c.manager.Do(ctx, func(ctx context.Context) error {
		err = c.storage.Create(ctx, game)
		if err != nil {
			return err
		}

		err = c.mediator.Publish(ctx, game.Events()...)
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

var _ core.CommandHandler = (*CreateGame)(nil)
