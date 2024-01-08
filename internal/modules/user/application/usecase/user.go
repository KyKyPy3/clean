package usecase

import (
	"context"
	"time"

	"github.com/avito-tech/go-transaction-manager/trm/v2/manager"

	"github.com/KyKyPy3/clean/internal/domain/common"
	"github.com/KyKyPy3/clean/internal/modules/user/domain/entity"
	"github.com/KyKyPy3/clean/internal/modules/user/infrastructure/controller/http/v1"
	"github.com/KyKyPy3/clean/pkg/logger"
)

const (
	requestTimeout = time.Second * 5
)

type UserService interface {
	Fetch(ctx context.Context, limit, offset int64) ([]entity.User, error)
	Create(ctx context.Context, data entity.User) (entity.User, error)
	GetByEmail(ctx context.Context, email common.Email) (entity.User, error)
	GetByID(ctx context.Context, id common.UID) (entity.User, error)
	Delete(ctx context.Context, id common.UID) error
}

type userUsecase struct {
	userService UserService
	logger      logger.Logger
	manager     *manager.Manager
}

func NewUserUsecase(userService UserService, manager *manager.Manager, logger logger.Logger) v1.UserUsecase {
	return &userUsecase{
		userService: userService,
		manager:     manager,
		logger:      logger,
	}
}

func (u *userUsecase) Fetch(ctx context.Context, limit, offset int64) ([]entity.User, error) {
	ctx, cancel := context.WithTimeout(ctx, requestTimeout)
	defer cancel()

	users, err := u.userService.Fetch(ctx, limit, offset)
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

func (u *userUsecase) GetByID(ctx context.Context, id common.UID) (entity.User, error) {
	ctx, cancel := context.WithTimeout(ctx, requestTimeout)
	defer cancel()

	user, err := u.userService.GetByID(ctx, id)
	if err != nil {
		return entity.User{}, err
	}

	return user, nil
}

func (u *userUsecase) GetByEmail(ctx context.Context, email common.Email) (entity.User, error) {
	ctx, cancel := context.WithTimeout(ctx, requestTimeout)
	defer cancel()

	user, err := u.userService.GetByEmail(ctx, email)
	if err != nil {
		return entity.User{}, err
	}

	return user, nil
}

func (u *userUsecase) Delete(ctx context.Context, id common.UID) error {
	ctx, cancel := context.WithTimeout(ctx, requestTimeout)
	defer cancel()

	err := u.userService.Delete(ctx, id)
	if err != nil {
		return err
	}

	return nil
}
