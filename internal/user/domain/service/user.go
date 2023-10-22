package service

import (
	"context"
	"github.com/KyKyPy3/clean/internal/common"
	"github.com/KyKyPy3/clean/internal/user/domain/entity"
	"github.com/KyKyPy3/clean/internal/user/usecase"
	"github.com/KyKyPy3/clean/pkg/logger"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
)

type UserPgStorage interface {
	Fetch(ctx context.Context, limit int64) ([]entity.User, error)
	Create(ctx context.Context, user entity.User) (entity.User, error)
	GetByID(ctx context.Context, id common.ID) (entity.User, error)
	GetByEmail(ctx context.Context, email string) (entity.User, error)
	Delete(ctx context.Context, id common.ID) error
}

type UserRedisStorage interface {
	GetByID(ctx context.Context, id uuid.UUID) (entity.User, error)
	Set(ctx context.Context, key string, user entity.User) error
	Delete(ctx context.Context, key string) error
}

type userService struct {
	pgStorage    UserPgStorage
	redisStorage UserRedisStorage
	logger       logger.Logger
}

func NewUserService(pgStorage UserPgStorage, redisStorage UserRedisStorage, logger logger.Logger) usecase.UserService {
	return &userService{pgStorage: pgStorage, redisStorage: redisStorage, logger: logger}
}

func (u *userService) Fetch(ctx context.Context, limit int64) ([]entity.User, error) {
	ctx, span := otel.Tracer("").Start(ctx, "userService.Fetch")
	defer span.End()

	users, err := u.pgStorage.Fetch(ctx, limit)
	if err != nil {
		return []entity.User{}, err
	}

	return users, nil
}

func (u *userService) Create(ctx context.Context, data entity.User) (entity.User, error) {
	ctx, span := otel.Tracer("").Start(ctx, "userService.Create")
	defer span.End()

	existed, err := u.pgStorage.GetByEmail(ctx, data.Email)
	if err != nil {
		return entity.User{}, err
	}

	if existed != (entity.User{}) {
		return entity.User{}, common.ErrEntityExist
	}

	user, err := u.pgStorage.Create(ctx, data)
	if err != nil {
		return entity.User{}, err
	}

	return user, nil
}

func (u *userService) GetByEmail(ctx context.Context, email string) (entity.User, error) {
	ctx, span := otel.Tracer("").Start(ctx, "userService.GetByEmail")
	defer span.End()

	user, err := u.pgStorage.GetByEmail(ctx, email)
	if err != nil {
		return entity.User{}, err
	}

	return user, nil
}

func (u *userService) GetByID(ctx context.Context, id common.ID) (entity.User, error) {
	ctx, span := otel.Tracer("").Start(ctx, "userService.GetByID")
	defer span.End()

	user, err := u.pgStorage.GetByID(ctx, id)
	if err != nil {
		return entity.User{}, err
	}

	return user, nil
}

func (u *userService) Delete(ctx context.Context, id common.ID) error {
	ctx, span := otel.Tracer("").Start(ctx, "userService.Delete")
	defer span.End()

	existed, err := u.pgStorage.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if existed == (entity.User{}) {
		return common.ErrNotFound
	}

	err = u.pgStorage.Delete(ctx, id)
	if err != nil {
		return err
	}

	return nil
}
