package ports

import (
	"context"
	"github.com/KyKyPy3/clean/internal/domain/common"
	"github.com/KyKyPy3/clean/internal/modules/user/domain/entity"
)

type UserPgStorage interface {
	Fetch(ctx context.Context, limit, offset int64) ([]entity.User, error)
	Create(ctx context.Context, data entity.User) error
	Update(ctx context.Context, data entity.User) error
	GetByEmail(ctx context.Context, email common.Email) (entity.User, error)
	GetByID(ctx context.Context, id common.UID) (entity.User, error)
	Delete(ctx context.Context, id common.UID) error
}

type UserRedisStorage interface {
	GetByID(ctx context.Context, id common.UID) (entity.User, error)
	Set(ctx context.Context, key string, user entity.User) error
	Delete(ctx context.Context, key string) error
}
