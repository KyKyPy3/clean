package entity_test

import (
	"testing"

	"github.com/KyKyPy3/clean/internal/common"
	"github.com/KyKyPy3/clean/internal/user/domain/entity"
	"github.com/stretchr/testify/assert"
)

func TestNewUser(t *testing.T) {
	user, err := entity.NewUser("Alise", "Cooper", "Lee", "alise@email.com", "12345")
	assert.Nil(t, err)
	assert.Equal(t, user.FirstName(), "Alise")
	assert.Equal(t, user.LastName(), "Cooper")
	assert.Equal(t, user.MiddleName(), "Lee")
	assert.Equal(t, user.Email(), "alise@email.com")
	assert.NotEqual(t, user.Password(), "12345")
}

func TestPasswordValidation(t *testing.T) {
	user, _ := entity.NewUser("Alise", "Cooper", "Lee", "alise@email.com", "12345")
	err := user.ValidatePassword("12345")
	assert.Nil(t, err)
	err = user.ValidatePassword("password")
	assert.NotNil(t, err)
}

func TestUserValidation(t *testing.T) {
	_, err := entity.NewUser("Alise", "Cooper", "Lee", "", "12345")
	assert.NotNil(t, err)
	assert.ErrorIs(t, err, common.ErrInvalidEntity)
}
