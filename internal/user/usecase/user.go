package usecase

import (
	"context"
	"github.com/KyKyPy3/clean/internal/common"
	v1 "github.com/KyKyPy3/clean/internal/user/controller/http/v1"
	"github.com/KyKyPy3/clean/internal/user/domain/entity"
	"github.com/KyKyPy3/clean/pkg/logger"
	"github.com/google/uuid"
	"time"
)

const (
	requestTimeout = time.Second * 5
)

type UserService interface {
	Fetch(ctx context.Context, limit int64) ([]entity.User, error)
	Create(ctx context.Context, data entity.User) (entity.User, error)
	GetByEmail(ctx context.Context, email string) (entity.User, error)
	GetByID(ctx context.Context, id common.ID) (entity.User, error)
	Delete(ctx context.Context, id common.ID) error
}

type userUsecase struct {
	userService UserService
	logger      logger.Logger
}

func NewUserUsecase(userService UserService, logger logger.Logger) v1.UserUsecase {
	return &userUsecase{
		userService: userService,
		logger:      logger,
	}
}

func (u *userUsecase) Fetch(ctx context.Context, limit int64) ([]entity.User, error) {
	ctx, cancel := context.WithTimeout(ctx, requestTimeout)
	defer cancel()

	users, err := u.userService.Fetch(ctx, limit)
	if err != nil {
		return []entity.User{}, err
	}

	return users, nil
}

func (u *userUsecase) Create(ctx context.Context, data entity.User) (entity.User, error) {
	ctx, cancel := context.WithTimeout(ctx, requestTimeout)
	defer cancel()

	user, err := u.userService.Create(ctx, data)
	if err != nil {
		return entity.User{}, err
	}

	return user, nil
}

func (u *userUsecase) GetByID(ctx context.Context, id uuid.UUID) (entity.User, error) {
	ctx, cancel := context.WithTimeout(ctx, requestTimeout)
	defer cancel()

	user, err := u.userService.GetByID(ctx, id)
	if err != nil {
		return entity.User{}, err
	}

	return user, nil
}

func (u *userUsecase) GetByEmail(ctx context.Context, email string) (entity.User, error) {
	ctx, cancel := context.WithTimeout(ctx, requestTimeout)
	defer cancel()

	user, err := u.userService.GetByEmail(ctx, email)
	if err != nil {
		return entity.User{}, err
	}

	return user, nil
}

func (u *userUsecase) Delete(ctx context.Context, id common.ID) error {
	ctx, cancel := context.WithTimeout(ctx, requestTimeout)
	defer cancel()

	err := u.userService.Delete(ctx, id)
	if err != nil {
		return err
	}

	return nil
}
