package entity

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/KyKyPy3/clean/internal/domain/common"
	"github.com/KyKyPy3/clean/internal/domain/core"
)

var (
	uniqueErr = errors.New("unique error")
)

type policyMock struct {
}

func (p *policyMock) IsUnique(email common.Email) (bool, error) {
	if email.String() == "not_unique@gmail.com" {
		return false, nil
	}

	if email.String() == "error@gmail.com" {
		return false, uniqueErr
	}

	return true, nil
}

func TestNewRegistration(t *testing.T) {
	email := common.MustNewEmail("alise@email.com")

	registration, err := NewRegistration(email, "", &policyMock{})
	assert.Nil(t, err)
	assert.Equal(t, registration.Email(), email)
	assert.Equal(t, registration.Verified(), false)
}

func TestRegistrationValidation(t *testing.T) {
	email := common.Email{}

	_, err := NewRegistration(email, "", &policyMock{})
	assert.NotNil(t, err)
	assert.ErrorIs(t, err, core.ErrInvalidEntity)
}

func TestRegistrationUniqueSuccess(t *testing.T) {
	email, _ := common.NewEmail("not_unique@gmail.com")

	_, err := NewRegistration(email, "", &policyMock{})
	assert.NotNil(t, err)
	assert.ErrorIs(t, err, core.ErrAlreadyExist)
}

func TestRegistrationUniqueError(t *testing.T) {
	email, _ := common.NewEmail("error@gmail.com")

	_, err := NewRegistration(email, "", &policyMock{})
	assert.NotNil(t, err)
	assert.ErrorIs(t, err, uniqueErr)
}
