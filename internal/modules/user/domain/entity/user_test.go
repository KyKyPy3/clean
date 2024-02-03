package entity

import (
	"github.com/KyKyPy3/clean/internal/domain/core"
	"golang.org/x/crypto/bcrypt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/KyKyPy3/clean/internal/domain/common"
	"github.com/KyKyPy3/clean/internal/modules/user/domain/value_object"
)

func TestNewUser(t *testing.T) {
	fullName := value_object.MustNewFullName("Alise", "Cooper", "Lee")
	email := common.MustNewEmail("alise@email.com")

	user, err := NewUser(fullName, email, "12345")
	assert.Nil(t, err)
	assert.Equal(t, user.FullName(), fullName)
	assert.Equal(t, user.Email(), email)
	assert.Equal(t, user.Password(), "12345")
}

func TestPasswordValidation(t *testing.T) {
	fullName := value_object.MustNewFullName("Alise", "Cooper", "Lee")
	email := common.MustNewEmail("alise@email.com")

	hash, _ := bcrypt.GenerateFromPassword([]byte("12345"), 10)

	user, _ := NewUser(fullName, email, string(hash))
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
	assert.ErrorIs(t, err, core.ErrInvalidEntity)
}
