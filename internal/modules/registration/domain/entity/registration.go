package entity

import (
	"fmt"

	"github.com/KyKyPy3/clean/internal/domain"
	"github.com/KyKyPy3/clean/internal/domain/common"
)

// Registration struct
type Registration struct {
	id common.UID

	email    common.Email
	verified bool
}

// NewRegistration - create and validate registration
func NewRegistration(email common.Email) (Registration, error) {
	if email.IsEmpty() {
		return Registration{}, fmt.Errorf("registration email is empty, err: %w", domain.ErrInvalidEntity)
	}

	r := Registration{
		email:    email,
		verified: false,
	}
	r.SetID(common.NewUID())

	return r, nil
}

func (r *Registration) GetID() common.UID {
	return r.id
}

func (r *Registration) SetID(id common.UID) {
	r.id = id
}

func (r *Registration) IsEmpty() bool {
	return *r == Registration{}
}

func (r *Registration) GetEmail() common.Email {
	return r.email
}

func (r *Registration) SetEmail(email common.Email) {
	r.email = email
}

func (r *Registration) GetVerified() bool {
	return r.verified
}

func (r *Registration) SetVerified(verified bool) {
	r.verified = verified
}
