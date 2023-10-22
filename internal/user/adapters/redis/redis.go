package redis

import (
	"context"
	"github.com/KyKyPy3/clean/internal/user/domain/entity"
	"github.com/KyKyPy3/clean/internal/user/domain/service"
	"github.com/KyKyPy3/clean/pkg/logger"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type userRedisStorage struct {
	db     *redis.Client
	logger logger.Logger
}

func NewUserRedisStorage(db *redis.Client, logger logger.Logger) service.UserRedisStorage {
	return &userRedisStorage{db: db, logger: logger}
}

func (u *userRedisStorage) GetByID(ctx context.Context, id uuid.UUID) (entity.User, error) {
	return entity.User{}, nil
}

func (u *userRedisStorage) Set(ctx context.Context, key string, user entity.User) error {
	return nil
}

func (u *userRedisStorage) Delete(ctx context.Context, key string) error {
	return nil
}
