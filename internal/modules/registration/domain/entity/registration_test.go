package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/KyKyPy3/clean/internal/domain"
	"github.com/KyKyPy3/clean/internal/domain/common"
)

func TestNewRegistration(t *testing.T) {
	email := common.MustNewEmail("alise@email.com")

	registration, err := NewRegistration(email)
	assert.Nil(t, err)
	assert.Equal(t, registration.GetEmail(), email)
	assert.Equal(t, registration.GetVerified(), false)
}

func TestRegistrationValidation(t *testing.T) {
	email := common.Email{}

	_, err := NewRegistration(email)
	assert.NotNil(t, err)
	assert.ErrorIs(t, err, domain.ErrInvalidEntity)
}
