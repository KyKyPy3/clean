package entity_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/KyKyPy3/clean/internal/domain/common"
	"github.com/KyKyPy3/clean/internal/domain/core"
	"github.com/KyKyPy3/clean/internal/modules/registration/domain/entity"
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

func TestNewRegistration(t *testing.T) {
	email := common.MustNewEmail("alise@email.com")

	registration, err := entity.NewRegistration(email, "", &policyMock{})
	require.NoError(t, err)
	assert.Equal(t, registration.Email(), email)
	assert.False(t, registration.Verified())
}

func TestRegistrationValidation(t *testing.T) {
	email := common.Email{}

	_, err := entity.NewRegistration(email, "", &policyMock{})
	require.Error(t, err)
	assert.ErrorIs(t, err, core.ErrInvalidEntity)
}

func TestRegistrationUniqueSuccess(t *testing.T) {
	email, _ := common.NewEmail("not_unique@gmail.com")

	_, err := entity.NewRegistration(email, "", &policyMock{})
	require.Error(t, err)
	assert.ErrorIs(t, err, core.ErrAlreadyExist)
}

func TestRegistrationUniqueError(t *testing.T) {
	email, _ := common.NewEmail("error@gmail.com")

	_, err := entity.NewRegistration(email, "", &policyMock{})
	require.Error(t, err)
	assert.ErrorIs(t, err, errUnique)
}
