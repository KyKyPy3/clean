package command

import (
	"context"
	"fmt"

	"github.com/KyKyPy3/clean/internal/application/core"
	"github.com/KyKyPy3/clean/internal/domain/common"
	"github.com/KyKyPy3/clean/internal/modules/session/application/ports"
	"github.com/KyKyPy3/clean/pkg/logger"
)

const LogoutUserKind = "LogoutUser"

type LogoutUserCommand struct {
	RefreshTokenID string
	AccessTokenID  string
}

func (c LogoutUserCommand) Type() core.CommandType {
	return LogoutUserKind
}

var _ core.Command = (*LogoutUserCommand)(nil)

type LogoutUser struct {
	sessionStorage ports.SessionRedisStorage
	logger         logger.Logger
}

func NewLogoutUser(
	sessionStorage ports.SessionRedisStorage,
	logger logger.Logger,
) LogoutUser {
	return LogoutUser{
		sessionStorage: sessionStorage,
		logger:         logger,
	}
}

func (l LogoutUser) Handle(ctx context.Context, command core.Command) (any, error) {
	logoutCommand, ok := command.(LogoutUserCommand)
	if !ok {
		return nil, fmt.Errorf("command type %s: %w", command.Type(), core.ErrUnexpectedCommand)
	}

	accessTokenID, err := common.ParseUID(logoutCommand.AccessTokenID)
	if err != nil {
		return nil, err
	}

	refreshTokenID, err := common.ParseUID(logoutCommand.RefreshTokenID)
	if err != nil {
		return nil, err
	}

	err = l.sessionStorage.Delete(ctx, accessTokenID)
	if err != nil {
		return nil, err
	}

	err = l.sessionStorage.Delete(ctx, refreshTokenID)
	if err != nil {
		return nil, err
	}

	return nil, nil
}
