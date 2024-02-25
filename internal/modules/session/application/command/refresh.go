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

const RefreshSessionKind = "RefreshSession"

type RefreshSessionCommand struct {
	ID         string
	UserID     string
	AccessTTL  time.Duration
	RefreshTTL time.Duration
}

type RefreshSessionResult struct {
	UserID      common.UID
	AccessToken entity.Token
}

func (c RefreshSessionCommand) Type() core.CommandType {
	return RefreshSessionKind
}

var _ core.Command = (*RefreshSessionCommand)(nil)

type RefreshSession struct {
	userView       ports.UserPgStorage
	sessionStorage ports.SessionRedisStorage
	logger         logger.Logger
}

func NewRefreshSession(
	userView ports.UserPgStorage,
	sessionStorage ports.SessionRedisStorage,
	logger logger.Logger,
) RefreshSession {
	return RefreshSession{
		userView:       userView,
		sessionStorage: sessionStorage,
		logger:         logger,
	}
}

func (l RefreshSession) Handle(ctx context.Context, command core.Command) (any, error) {
	refreshCommand, ok := command.(RefreshSessionCommand)
	if !ok {
		return nil, fmt.Errorf("command type %s: %w", command.Type(), core.ErrUnexpectedCommand)
	}

	userID, err := common.ParseUID(refreshCommand.UserID)
	if err != nil {
		return nil, domain_core.ErrNotFound
	}

	user, err := l.userView.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	if user.IsEmpty() {
		return nil, domain_core.ErrNotFound
	}

	now := time.Now().UTC()
	accessExpiresIn := now.Add(refreshCommand.AccessTTL).Unix()
	accessToken := entity.NewToken(user.ID(), accessExpiresIn)

	err = l.sessionStorage.Set(ctx, accessToken.ID(), accessToken)
	if err != nil {
		return nil, err
	}

	return RefreshSessionResult{
		UserID:      user.ID(),
		AccessToken: accessToken,
	}, nil
}
