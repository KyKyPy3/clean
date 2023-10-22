package service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/KyKyPy3/clean/config"
	"github.com/KyKyPy3/clean/internal/common"
	"github.com/KyKyPy3/clean/internal/user/domain/entity"
	"github.com/KyKyPy3/clean/internal/user/domain/service"
	"github.com/KyKyPy3/clean/internal/user/usecase"
	mocks "github.com/KyKyPy3/clean/mocks/internal_/user/domain/service"
	"github.com/KyKyPy3/clean/pkg/logger"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func prepare(t *testing.T) (usecase.UserService, *mocks.UserPgStorage, *mocks.UserRedisStorage) {
	t.Helper()

	mockPgUserStorage := new(mocks.UserPgStorage)
	mockRedisuserStorage := new(mocks.UserRedisStorage)

	loggerCfg := &config.LoggerConfig{
		Mode:     "development",
		Level:    "debug",
		Encoding: "json",
	}

	logger := logger.NewLogger(loggerCfg)

	service := service.NewUserService(mockPgUserStorage, mockRedisuserStorage, logger)

	return service, mockPgUserStorage, mockRedisuserStorage
}

func TestFetch(t *testing.T) {
	t.Parallel()

	service, mockPgUserStorage, _ := prepare(t)

	mockUser := entity.User{
		Name:       "Alise",
		Surname:    "Cooper",
		Middlename: "Saint",
		Email:      "alise@email.com",
	}
	mockUsersList := make([]entity.User, 0)
	mockUsersList = append(mockUsersList, mockUser)

	t.Run("success", func(t *testing.T) {
		mockPgUserStorage.EXPECT().Fetch(mock.Anything, mock.AnythingOfType("int64")).Return(mockUsersList, nil).Once()
		list, err := service.Fetch(context.Background(), int64(1))
		assert.NoError(t, err)
		assert.Equal(t, mockUsersList, list)

		mockPgUserStorage.AssertExpectations(t)
	})

	t.Run("error", func(t *testing.T) {
		mockPgUserStorage.EXPECT().Fetch(mock.Anything, mock.AnythingOfType("int64")).Return(nil, errors.New("Unexpected")).Once()
		list, err := service.Fetch(context.TODO(), int64(1))

		assert.Error(t, err)
		assert.Len(t, list, 0)

		mockPgUserStorage.AssertExpectations(t)
	})
}

func TestGetByID(t *testing.T) {
	t.Parallel()

	service, mockPgUserStorage, _ := prepare(t)

	mockUser := entity.User{
		ID:         uuid.New(),
		Name:       "Alise",
		Surname:    "Cooper",
		Middlename: "Saint",
		Email:      "alise@email.com",
	}

	t.Run("success", func(t *testing.T) {
		mockPgUserStorage.EXPECT().GetByID(mock.Anything, mock.Anything).Return(mockUser, nil).Once()
		user, err := service.GetByID(context.Background(), mockUser.ID)
		assert.NoError(t, err)
		assert.Equal(t, mockUser, user)

		mockPgUserStorage.AssertExpectations(t)
	})

	t.Run("error", func(t *testing.T) {
		mockPgUserStorage.EXPECT().GetByID(mock.Anything, mock.Anything).Return(entity.User{}, errors.New("Unexpected")).Once()
		user, err := service.GetByID(context.Background(), mockUser.ID)

		assert.Error(t, err)
		assert.Equal(t, entity.User{}, user)

		mockPgUserStorage.AssertExpectations(t)
	})
}

func TestGetByEmail(t *testing.T) {
	t.Parallel()

	service, mockPgUserStorage, _ := prepare(t)

	mockUser := entity.User{
		Name:       "Alise",
		Surname:    "Cooper",
		Middlename: "Saint",
		Email:      "alise@email.com",
	}

	t.Run("success", func(t *testing.T) {
		mockPgUserStorage.EXPECT().GetByEmail(mock.Anything, mock.AnythingOfType("string")).Return(mockUser, nil).Once()
		user, err := service.GetByEmail(context.Background(), mockUser.Email)
		assert.NoError(t, err)
		assert.Equal(t, mockUser, user)

		mockPgUserStorage.AssertExpectations(t)
	})

	t.Run("error", func(t *testing.T) {
		mockPgUserStorage.EXPECT().GetByEmail(mock.Anything, mock.AnythingOfType("string")).Return(entity.User{}, errors.New("Unexpected")).Once()
		user, err := service.GetByEmail(context.Background(), mockUser.Email)

		assert.Error(t, err)
		assert.Equal(t, entity.User{}, user)

		mockPgUserStorage.AssertExpectations(t)
	})
}

func TestCreate(t *testing.T) {
	t.Parallel()

	service, mockPgUserStorage, _ := prepare(t)

	mockUser := entity.User{
		Name:       "Alise",
		Surname:    "Cooper",
		Middlename: "Saint",
		Email:      "alise@email.com",
	}

	t.Run("success", func(t *testing.T) {
		tempMockUser := mockUser
		tempMockUser.ID = uuid.New()

		mockPgUserStorage.EXPECT().GetByEmail(mock.Anything, mock.AnythingOfType("string")).Return(entity.User{}, nil).Once()
		mockPgUserStorage.EXPECT().Create(mock.Anything, mock.AnythingOfType("entity.User")).Return(mockUser, nil).Once()

		user, err := service.Create(context.Background(), tempMockUser)

		assert.NoError(t, err)
		assert.Equal(t, mockUser, user)

		mockPgUserStorage.AssertExpectations(t)
	})

	t.Run("user exist", func(t *testing.T) {
		tempMockUser := mockUser
		tempMockUser.ID = uuid.New()

		mockPgUserStorage.EXPECT().GetByEmail(mock.Anything, mock.AnythingOfType("string")).Return(mockUser, nil).Once()

		user, err := service.Create(context.Background(), tempMockUser)

		assert.Error(t, err)
		assert.Equal(t, entity.User{}, user)
		assert.Equal(t, err, common.ErrEntityExist)

		mockPgUserStorage.AssertExpectations(t)
	})

	t.Run("error of user locate", func(t *testing.T) {
		tempMockUser := mockUser
		tempMockUser.ID = uuid.New()

		mockPgUserStorage.EXPECT().GetByEmail(mock.Anything, mock.AnythingOfType("string")).Return(entity.User{}, errors.New("Unexpected")).Once()

		user, err := service.Create(context.Background(), mockUser)

		assert.Error(t, err)
		assert.Equal(t, entity.User{}, user)

		mockPgUserStorage.AssertExpectations(t)
	})

	t.Run("error", func(t *testing.T) {
		tempMockUser := mockUser
		tempMockUser.ID = uuid.New()

		mockPgUserStorage.EXPECT().GetByEmail(mock.Anything, mock.AnythingOfType("string")).Return(entity.User{}, nil).Once()
		mockPgUserStorage.EXPECT().Create(mock.Anything, mock.AnythingOfType("entity.User")).Return(entity.User{}, errors.New("Unexpected")).Once()

		user, err := service.Create(context.Background(), mockUser)

		assert.Error(t, err)
		assert.Equal(t, entity.User{}, user)

		mockPgUserStorage.AssertExpectations(t)
	})
}

func TestDelete(t *testing.T) {
	t.Parallel()

	service, mockPgUserStorage, _ := prepare(t)

	mockUser := entity.User{
		Name:       "Alise",
		Surname:    "Cooper",
		Middlename: "Saint",
		Email:      "alise@email.com",
	}

	t.Run("success", func(t *testing.T) {
		mockPgUserStorage.EXPECT().GetByID(mock.Anything, mock.Anything).Return(mockUser, nil).Once()
		mockPgUserStorage.EXPECT().Delete(mock.Anything, mock.Anything).Return(nil).Once()

		err := service.Delete(context.Background(), mockUser.ID)
		assert.NoError(t, err)

		mockPgUserStorage.AssertExpectations(t)
	})

	t.Run("user not exist", func(t *testing.T) {
		mockPgUserStorage.EXPECT().GetByID(mock.Anything, mock.Anything).Return(entity.User{}, nil).Once()

		err := service.Delete(context.Background(), mockUser.ID)

		assert.Error(t, err)

		mockPgUserStorage.AssertExpectations(t)
	})

	t.Run("error check user", func(t *testing.T) {
		mockPgUserStorage.EXPECT().GetByID(mock.Anything, mock.Anything).Return(entity.User{}, errors.New("Unexpected")).Once()

		err := service.Delete(context.Background(), mockUser.ID)

		assert.Error(t, err)

		mockPgUserStorage.AssertExpectations(t)
	})

	t.Run("error", func(t *testing.T) {
		mockPgUserStorage.EXPECT().GetByID(mock.Anything, mock.Anything).Return(mockUser, nil).Once()
		mockPgUserStorage.EXPECT().Delete(mock.Anything, mock.Anything).Return(errors.New("Unexpected")).Once()

		err := service.Delete(context.Background(), mockUser.ID)

		assert.Error(t, err)

		mockPgUserStorage.AssertExpectations(t)
	})
}
