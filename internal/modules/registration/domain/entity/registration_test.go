package entity

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/KyKyPy3/clean/internal/domain/common"
	"github.com/KyKyPy3/clean/internal/domain/core"
)

type policyMock struct {
}

func (p *policyMock) IsUnique(context.Context, common.Email) (bool, error) {
	return true, nil
}

func TestNewRegistration(t *testing.T) {
	email := common.MustNewEmail("alise@email.com")

	registration, err := NewRegistration(context.Background(), email, &policyMock{})
	assert.Nil(t, err)
	assert.Equal(t, registration.Email(), email)
	assert.Equal(t, registration.Verified(), false)
}

func TestRegistrationValidation(t *testing.T) {
	email := common.Email{}

	_, err := NewRegistration(context.Background(), email, &policyMock{})
	assert.NotNil(t, err)
	assert.ErrorIs(t, err, core.ErrInvalidEntity)
}
