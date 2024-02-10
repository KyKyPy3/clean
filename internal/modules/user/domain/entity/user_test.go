package entity_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"

	"github.com/KyKyPy3/clean/internal/domain/common"
	"github.com/KyKyPy3/clean/internal/domain/core"
	"github.com/KyKyPy3/clean/internal/modules/user/domain/entity"
	"github.com/KyKyPy3/clean/internal/modules/user/domain/vo"
)

var (
	errUnique = errors.New("unique error")
)

type policyMock struct {
}

func (p *policyMock) IsUnique(email common.Email) (bool, error) {
	if email.String() == "not_unique@gmail.com" {
		return false, nil
	}

	if email.String() == "error@gmail.com" {
		return false, errUnique
	}

	return true, nil
}

func TestNewUser(t *testing.T) {
	fullName := vo.MustNewFullName("Alise", "Cooper", "Lee")
	email := common.MustNewEmail("alise@email.com")

	user, err := entity.NewUser(fullName, email, "12345", &policyMock{})
	require.NoError(t, err)
	assert.Equal(t, user.FullName(), fullName)
	assert.Equal(t, user.Email(), email)
	assert.Equal(t, "12345", user.Password())
}

func TestPasswordValidation(t *testing.T) {
	fullName := vo.MustNewFullName("Alise", "Cooper", "Lee")
	email := common.MustNewEmail("alise@email.com")

	hash, _ := bcrypt.GenerateFromPassword([]byte("12345"), 10)

	user, _ := entity.NewUser(fullName, email, string(hash), &policyMock{})
	err := user.ValidatePassword("12345")
	require.NoError(t, err)
	err = user.ValidatePassword("password")
	assert.Error(t, err)
}

func TestUserValidation(t *testing.T) {
	fullName := vo.MustNewFullName("Alise", "Cooper", "Lee")
	email := common.Email{}

	_, err := entity.NewUser(fullName, email, "12345", &policyMock{})
	require.Error(t, err)
	assert.ErrorIs(t, err, core.ErrInvalidEntity)
}
