package redis

import (
	"context"

	"github.com/redis/go-redis/v9"
	"go.opentelemetry.io/otel"

	"github.com/KyKyPy3/clean/internal/domain/common"
	"github.com/KyKyPy3/clean/internal/modules/user/application/ports"
	"github.com/KyKyPy3/clean/internal/modules/user/domain/entity"
	"github.com/KyKyPy3/clean/pkg/logger"
)

type userRedisStorage struct {
	db     *redis.Client
	logger logger.Logger
}

func NewUserRedisStorage(db *redis.Client, logger logger.Logger) ports.UserRedisStorage {
	return &userRedisStorage{db: db, logger: logger}
}

func (u *userRedisStorage) GetByID(ctx context.Context, id common.UID) (entity.User, error) {
	_, span := otel.Tracer("").Start(ctx, "userRedisStorage.GetByID")
	defer span.End()

	u.logger.Debugf("ID: %s", id)

	return entity.User{}, nil
}

func (u *userRedisStorage) Set(ctx context.Context, key string, user entity.User) error {
	_, span := otel.Tracer("").Start(ctx, "userRedisStorage.Set")
	defer span.End()

	u.logger.Debugf("Key: %s", key)
	u.logger.Debugf("User: %v", user)

	return nil
}

func (u *userRedisStorage) Delete(ctx context.Context, key string) error {
	_, span := otel.Tracer("").Start(ctx, "userRedisStorage.Delete")
	defer span.End()

	u.logger.Debugf("Key: %s", key)

	return nil
}
