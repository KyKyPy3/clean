package ports

import (
	"context"

	"github.com/KyKyPy3/clean/internal/domain/common"
	"github.com/KyKyPy3/clean/internal/modules/session/domain/entity"
	user_domain "github.com/KyKyPy3/clean/internal/modules/user/domain/entity"
)

type UserPgStorage interface {
	Fetch(ctx context.Context, limit, offset int64) ([]user_domain.User, error)
	GetByEmail(ctx context.Context, email common.Email) (user_domain.User, error)
	GetByID(ctx context.Context, id common.UID) (user_domain.User, error)
}

type SessionRedisStorage interface {
	Get(ctx context.Context, tokenID common.UID) (entity.Token, error)
	Set(ctx context.Context, tokenID common.UID, token entity.Token) error
	Delete(ctx context.Context, tokenID common.UID) error
}
