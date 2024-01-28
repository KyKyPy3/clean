package service

import (
	"context"
	"errors"
	"github.com/KyKyPy3/clean/internal/domain/core"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"

	"github.com/KyKyPy3/clean/internal/domain/common"
	"github.com/KyKyPy3/clean/internal/modules/user/application/usecase"
	"github.com/KyKyPy3/clean/internal/modules/user/domain/entity"
	"github.com/KyKyPy3/clean/pkg/logger"
)

type UserPgStorage interface {
	Fetch(ctx context.Context, limit, offset int64) ([]entity.User, error)
	Create(ctx context.Context, user entity.User) (entity.User, error)
	GetByID(ctx context.Context, id common.UID) (entity.User, error)
	GetByEmail(ctx context.Context, email common.Email) (entity.User, error)
	Delete(ctx context.Context, id common.UID) error
}

type UserRedisStorage interface {
	GetByID(ctx context.Context, id common.UID) (entity.User, error)
	Set(ctx context.Context, key string, user entity.User) error
	Delete(ctx context.Context, key string) error
}

type userService struct {
	pgStorage    UserPgStorage
	redisStorage UserRedisStorage
	logger       logger.Logger
	tracer       trace.Tracer
}

func NewUserService(pgStorage UserPgStorage, redisStorage UserRedisStorage, logger logger.Logger) usecase.UserService {
	return &userService{
		pgStorage:    pgStorage,
		redisStorage: redisStorage,
		logger:       logger,
		tracer:       otel.Tracer(""),
	}
}

func (u *userService) Fetch(ctx context.Context, limit, offset int64) ([]entity.User, error) {
	ctx, span := u.tracer.Start(ctx, "userService.Fetch")
	defer span.End()

	users, err := u.pgStorage.Fetch(ctx, limit, offset)
	if err != nil {
		return []entity.User{}, err
	}

	return users, nil
}

func (u *userService) Create(ctx context.Context, data entity.User) (entity.User, error) {
	ctx, span := u.tracer.Start(ctx, "userService.Create")
	defer span.End()

	existed, err := u.pgStorage.GetByEmail(ctx, data.GetEmail())
	if err != nil && !errors.Is(err, core.ErrNotFound) {
		return entity.User{}, err
	}

	if !existed.IsEmpty() {
		return entity.User{}, core.ErrAlreadyExist
	}

	user, err := u.pgStorage.Create(ctx, data)
	if err != nil {
		return entity.User{}, err
	}

	return user, nil
}

func (u *userService) GetByEmail(ctx context.Context, email common.Email) (entity.User, error) {
	ctx, span := u.tracer.Start(ctx, "userService.GetByEmail")
	defer span.End()

	user, err := u.pgStorage.GetByEmail(ctx, email)
	if err != nil {
		return entity.User{}, err
	}

	return user, nil
}

func (u *userService) GetByID(ctx context.Context, id common.UID) (entity.User, error) {
	ctx, span := u.tracer.Start(ctx, "userService.GetByID")
	defer span.End()

	user, err := u.pgStorage.GetByID(ctx, id)
	if err != nil {
		return entity.User{}, err
	}

	return user, nil
}

func (u *userService) Delete(ctx context.Context, id common.UID) error {
	ctx, span := u.tracer.Start(ctx, "userService.Delete")
	defer span.End()

	existed, err := u.pgStorage.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if existed.IsEmpty() {
		return core.ErrNotFound
	}

	err = u.pgStorage.Delete(ctx, id)
	if err != nil {
		return err
	}

	return nil
}
