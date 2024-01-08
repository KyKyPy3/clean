package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/KyKyPy3/clean/internal/domain"
	"github.com/KyKyPy3/clean/internal/domain/common"
	"github.com/KyKyPy3/clean/internal/modules/user/domain/value_object"
)

func TestNewUser(t *testing.T) {
	fullName := value_object.MustNewFullName("Alise", "Cooper", "Lee")
	email := common.MustNewEmail("alise@email.com")

	user, err := NewUser(fullName, email, "12345")
	assert.Nil(t, err)
	assert.Equal(t, user.GetFullName(), fullName)
	assert.Equal(t, user.GetEmail(), email)
	assert.NotEqual(t, user.GetPassword(), "12345")
}

func TestPasswordValidation(t *testing.T) {
	fullName := value_object.MustNewFullName("Alise", "Cooper", "Lee")
	email := common.MustNewEmail("alise@email.com")

	user, _ := NewUser(fullName, email, "12345")
	err := user.ValidatePassword("12345")
	assert.Nil(t, err)
	err = user.ValidatePassword("password")
	assert.NotNil(t, err)
}

func TestUserValidation(t *testing.T) {
	fullName := value_object.MustNewFullName("Alise", "Cooper", "Lee")
	email := common.Email{}

	_, err := NewUser(fullName, email, "12345")
	assert.NotNil(t, err)
	assert.ErrorIs(t, err, domain.ErrInvalidEntity)
}
