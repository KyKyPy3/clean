package command

import (
	"context"
	"fmt"
	"time"

	"github.com/KyKyPy3/clean/internal/application/core"
	"github.com/KyKyPy3/clean/internal/domain/common"
	domain_core "github.com/KyKyPy3/clean/internal/domain/core"
	"github.com/KyKyPy3/clean/internal/modules/session/application/ports"
	"github.com/KyKyPy3/clean/internal/modules/session/domain/entity"
	"github.com/KyKyPy3/clean/pkg/logger"
)

const LoginUserKind = "LoginUser"

type LoginUserCommand struct {
	Email      string
	Password   string
	AccessTTL  time.Duration
	RefreshTTL time.Duration
}

type LoginUserResult struct {
	UserID       common.UID
	AccessToken  entity.Token
	RefreshToken entity.Token
}

func (c LoginUserCommand) Type() core.CommandType {
	return LoginUserKind
}

var _ core.Command = (*LoginUserCommand)(nil)

type LoginUser struct {
	userView       ports.UserPgStorage
	sessionStorage ports.SessionRedisStorage
	logger         logger.Logger
}

func NewLoginUser(
	userView ports.UserPgStorage,
	sessionStorage ports.SessionRedisStorage,
	logger logger.Logger,
) LoginUser {
	return LoginUser{
		userView:       userView,
		sessionStorage: sessionStorage,
		logger:         logger,
	}
}

func (l LoginUser) Handle(ctx context.Context, command core.Command) (any, error) {
	loginCommand, ok := command.(LoginUserCommand)
	if !ok {
		return nil, fmt.Errorf("command type %s: %w", command.Type(), core.ErrUnexpectedCommand)
	}

	email, err := common.NewEmail(loginCommand.Email)
	if err != nil {
		return nil, err
	}

	user, err := l.userView.GetByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	if user.IsEmpty() {
		return nil, domain_core.ErrNotFound
	}

	err = user.ValidatePassword(loginCommand.Password)
	if err != nil {
		return nil, domain_core.ErrNotFound
	}

	now := time.Now().UTC()
	accessExpiresIn := now.Add(loginCommand.AccessTTL).Unix()
	refreshExpiresIn := now.Add(loginCommand.RefreshTTL).Unix()
	accessToken := entity.NewToken(user.ID(), accessExpiresIn)
	refreshToken := entity.NewToken(user.ID(), refreshExpiresIn)

	err = l.sessionStorage.Set(ctx, accessToken.ID(), accessToken)
	if err != nil {
		return nil, err
	}

	err = l.sessionStorage.Set(ctx, refreshToken.ID(), refreshToken)
	if err != nil {
		return nil, err
	}

	return LoginUserResult{
		UserID:       user.ID(),
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
