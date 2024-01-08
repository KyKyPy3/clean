package service_test

import (
	"context"
	"errors"
	"github.com/KyKyPy3/clean/internal/domain/common"
	"github.com/KyKyPy3/clean/internal/infrastructure/config"
	"github.com/KyKyPy3/clean/internal/modules/user/application/usecase"
	"github.com/KyKyPy3/clean/internal/modules/user/domain/value_object"
	"testing"

	"github.com/KyKyPy3/clean/internal/domain"
	"github.com/KyKyPy3/clean/internal/modules/user/domain/entity"
	"github.com/KyKyPy3/clean/internal/modules/user/domain/service"
	mocks "github.com/KyKyPy3/clean/mocks/internal_/modules/user/domain/service"
	"github.com/KyKyPy3/clean/pkg/logger"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func prepare(t *testing.T) (usecase.UserService, *mocks.UserPgStorage, *mocks.UserRedisStorage) {
	t.Helper()

	mockPgUserStorage := new(mocks.UserPgStorage)
	mockRedisUserStorage := new(mocks.UserRedisStorage)

	loggerCfg := &config.LoggerConfig{
		Mode:     "development",
		Level:    "debug",
		Encoding: "json",
	}

	log := logger.NewLogger(loggerCfg)

	srv := service.NewUserService(mockPgUserStorage, mockRedisUserStorage, log)

	return srv, mockPgUserStorage, mockRedisUserStorage
}

func TestFetch(t *testing.T) {
	t.Parallel()

	srv, mockPgUserStorage, _ := prepare(t)
	fullName := value_object.MustNewFullName("Alise", "Cooper", "Saint")
	email := common.MustNewEmail("alise@email.com")

	mockUser := entity.User{}
	mockUser.SetFullName(fullName)
	mockUser.SetEmail(email)

	mockUsersList := make([]entity.User, 0)
	mockUsersList = append(mockUsersList, mockUser)

	t.Run("success", func(t *testing.T) {
		mockPgUserStorage.EXPECT().Fetch(
			mock.Anything,
			mock.AnythingOfType("int64"),
			mock.AnythingOfType("int64"),
		).Return(mockUsersList, nil).Once()
		list, err := srv.Fetch(context.Background(), int64(1), int64(0))
		assert.NoError(t, err)
		assert.Equal(t, mockUsersList, list)

		mockPgUserStorage.AssertExpectations(t)
	})

	t.Run("error", func(t *testing.T) {
		mockPgUserStorage.EXPECT().Fetch(
			mock.Anything,
			mock.AnythingOfType("int64"),
			mock.AnythingOfType("int64"),
		).Return(nil, errors.New("unexpected")).Once()
		list, err := srv.Fetch(context.TODO(), int64(1), int64(0))

		assert.Error(t, err)
		assert.Len(t, list, 0)

		mockPgUserStorage.AssertExpectations(t)
	})
}

func TestGetByID(t *testing.T) {
	t.Parallel()

	srv, mockPgUserStorage, _ := prepare(t)
	fullName := value_object.MustNewFullName("Alise", "Cooper", "Saint")
	email := common.MustNewEmail("alise@email.com")

	mockUser := entity.User{}
	mockUser.SetID(common.NewUID())
	mockUser.SetFullName(fullName)
	mockUser.SetEmail(email)

	t.Run("success", func(t *testing.T) {
		mockPgUserStorage.EXPECT().GetByID(mock.Anything, mock.Anything).Return(mockUser, nil).Once()
		user, err := srv.GetByID(context.Background(), mockUser.GetID())
		assert.NoError(t, err)
		assert.Equal(t, mockUser, user)

		mockPgUserStorage.AssertExpectations(t)
	})

	t.Run("error", func(t *testing.T) {
		mockPgUserStorage.EXPECT().GetByID(
			mock.Anything,
			mock.Anything,
		).Return(entity.User{}, errors.New("unexpected")).Once()
		user, err := srv.GetByID(context.Background(), mockUser.GetID())

		assert.Error(t, err)
		assert.Equal(t, entity.User{}, user)

		mockPgUserStorage.AssertExpectations(t)
	})
}

func TestGetByEmail(t *testing.T) {
	t.Parallel()

	srv, mockPgUserStorage, _ := prepare(t)
	fullName := value_object.MustNewFullName("Alise", "Cooper", "Saint")
	email := common.MustNewEmail("alise@email.com")

	mockUser := entity.User{}
	mockUser.SetFullName(fullName)
	mockUser.SetEmail(email)

	t.Run("success", func(t *testing.T) {
		mockPgUserStorage.EXPECT().GetByEmail(
			mock.Anything,
			mock.AnythingOfType("common.Email"),
		).Return(mockUser, nil).Once()
		user, err := srv.GetByEmail(context.Background(), mockUser.GetEmail())
		assert.NoError(t, err)
		assert.Equal(t, mockUser, user)

		mockPgUserStorage.AssertExpectations(t)
	})

	t.Run("error", func(t *testing.T) {
		mockPgUserStorage.EXPECT().GetByEmail(
			mock.Anything,
			mock.AnythingOfType("common.Email"),
		).Return(entity.User{}, errors.New("unexpected")).Once()
		user, err := srv.GetByEmail(context.Background(), mockUser.GetEmail())

		assert.Error(t, err)
		assert.Equal(t, entity.User{}, user)

		mockPgUserStorage.AssertExpectations(t)
	})
}

func TestCreate(t *testing.T) {
	t.Parallel()

	srv, mockPgUserStorage, _ := prepare(t)
	fullName := value_object.MustNewFullName("Alise", "Cooper", "Saint")
	email := common.MustNewEmail("alise@email.com")

	mockUser := entity.User{}
	mockUser.SetFullName(fullName)
	mockUser.SetEmail(email)

	t.Run("success", func(t *testing.T) {
		tempMockUser := mockUser
		tempMockUser.SetID(common.NewUID())

		mockPgUserStorage.EXPECT().GetByEmail(mock.Anything, mock.AnythingOfType("common.Email")).Return(entity.User{}, nil).Once()
		mockPgUserStorage.EXPECT().Create(mock.Anything, mock.AnythingOfType("entity.User")).Return(mockUser, nil).Once()

		user, err := srv.Create(context.Background(), tempMockUser)

		assert.NoError(t, err)
		assert.Equal(t, mockUser, user)

		mockPgUserStorage.AssertExpectations(t)
	})

	t.Run("user exist", func(t *testing.T) {
		tempMockUser := mockUser
		tempMockUser.SetID(common.NewUID())

		mockPgUserStorage.EXPECT().GetByEmail(mock.Anything, mock.AnythingOfType("common.Email")).Return(tempMockUser, nil).Once()

		user, err := srv.Create(context.Background(), tempMockUser)

		assert.Error(t, err)
		assert.Equal(t, entity.User{}, user)
		assert.Equal(t, err, domain.ErrAlreadyExist)

		mockPgUserStorage.AssertExpectations(t)
	})

	t.Run("error of user locate", func(t *testing.T) {
		tempMockUser := mockUser
		tempMockUser.SetID(common.NewUID())

		mockPgUserStorage.EXPECT().GetByEmail(mock.Anything, mock.AnythingOfType("common.Email")).Return(entity.User{}, errors.New("Unexpected")).Once()

		user, err := srv.Create(context.Background(), mockUser)

		assert.Error(t, err)
		assert.Equal(t, entity.User{}, user)

		mockPgUserStorage.AssertExpectations(t)
	})

	t.Run("error", func(t *testing.T) {
		tempMockUser := mockUser
		tempMockUser.SetID(common.NewUID())

		mockPgUserStorage.EXPECT().GetByEmail(mock.Anything, mock.AnythingOfType("common.Email")).Return(entity.User{}, nil).Once()
		mockPgUserStorage.EXPECT().Create(mock.Anything, mock.AnythingOfType("entity.User")).Return(entity.User{}, errors.New("Unexpected")).Once()

		user, err := srv.Create(context.Background(), mockUser)

		assert.Error(t, err)
		assert.Equal(t, entity.User{}, user)

		mockPgUserStorage.AssertExpectations(t)
	})
}

func TestDelete(t *testing.T) {
	t.Parallel()

	srv, mockPgUserStorage, _ := prepare(t)
	fullName := value_object.MustNewFullName("Alise", "Cooper", "Saint")
	email := common.MustNewEmail("alise@email.com")

	mockUser := entity.User{}
	mockUser.SetID(common.NewUID())
	mockUser.SetFullName(fullName)
	mockUser.SetEmail(email)

	t.Run("success", func(t *testing.T) {
		mockPgUserStorage.EXPECT().GetByID(mock.Anything, mock.Anything).Return(mockUser, nil).Once()
		mockPgUserStorage.EXPECT().Delete(mock.Anything, mock.Anything).Return(nil).Once()

		err := srv.Delete(context.Background(), mockUser.GetID())
		assert.NoError(t, err)

		mockPgUserStorage.AssertExpectations(t)
	})

	t.Run("user not exist", func(t *testing.T) {
		mockPgUserStorage.EXPECT().GetByID(mock.Anything, mock.Anything).Return(entity.User{}, nil).Once()

		err := srv.Delete(context.Background(), mockUser.GetID())

		assert.Error(t, err)

		mockPgUserStorage.AssertExpectations(t)
	})

	t.Run("error check user", func(t *testing.T) {
		mockPgUserStorage.EXPECT().GetByID(mock.Anything, mock.Anything).Return(entity.User{}, errors.New("Unexpected")).Once()

		err := srv.Delete(context.Background(), mockUser.GetID())

		assert.Error(t, err)

		mockPgUserStorage.AssertExpectations(t)
	})

	t.Run("error", func(t *testing.T) {
		mockPgUserStorage.EXPECT().GetByID(mock.Anything, mock.Anything).Return(mockUser, nil).Once()
		mockPgUserStorage.EXPECT().Delete(mock.Anything, mock.Anything).Return(errors.New("unexpected")).Once()

		err := srv.Delete(context.Background(), mockUser.GetID())

		assert.Error(t, err)

		mockPgUserStorage.AssertExpectations(t)
	})
}
