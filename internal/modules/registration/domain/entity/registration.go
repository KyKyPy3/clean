package entity

import (
	"context"
	"fmt"

	"golang.org/x/crypto/bcrypt"

	"github.com/KyKyPy3/clean/internal/domain/common"
	"github.com/KyKyPy3/clean/internal/domain/core"
	"github.com/KyKyPy3/clean/internal/modules/registration/domain"
	"github.com/KyKyPy3/clean/internal/modules/registration/domain/event"
	"github.com/KyKyPy3/clean/pkg/mediator"
)

// Registration struct
type Registration struct {
	*core.BaseAggregateRoot

	id       common.UID
	email    common.Email
	password string
	verified bool
}

// NewRegistration - create and validate registration
func NewRegistration(ctx context.Context, email common.Email, password string, uniqPolicy domain.UniqueEmailPolicy) (Registration, error) {
	if email.IsEmpty() {
		return Registration{}, fmt.Errorf("registration email is empty, err: %w", core.ErrInvalidEntity)
	}

	ok, err := uniqPolicy.IsUnique(ctx, email)
	if err != nil {
		return Registration{}, fmt.Errorf("failed to check uniqueness of email on registration, err: %w", err)
	}

	if !ok {
		return Registration{}, core.ErrAlreadyExist
	}

	r := Registration{
		BaseAggregateRoot: &core.BaseAggregateRoot{},
		id:                common.NewUID(),
		email:             email,
		password:          password,
		verified:          false,
	}

	if err := r.hashPassword(); err != nil {
		return Registration{}, err
	}

	r.BaseAggregateRoot.AddEvent(event.RegistrationCreatedEvent{ID: r.ID().String(), Email: email})

	return r, nil
}

func Hydrate(id common.UID, email common.Email, password string, verified bool) Registration {
	reg := Registration{
		BaseAggregateRoot: &core.BaseAggregateRoot{},
		id:                id,
		email:             email,
		password:          password,
		verified:          verified,
	}

	return reg
}

func (r *Registration) IsEmpty() bool {
	return *r == Registration{}
}

func (r *Registration) ID() common.UID {
	return r.id
}

func (r *Registration) Email() common.Email {
	return r.email
}

func (r *Registration) Password() string {
	return r.password
}

func (r *Registration) Verified() bool {
	return r.verified
}

func (r *Registration) Verify() error {
	if r.verified {
		return core.ErrNoChanges
	}

	r.verified = true
	r.BaseAggregateRoot.AddEvent(event.RegistrationVerifiedEvent{ID: r.id.String(), Email: r.email, Password: r.password})

	return nil
}

// hashPassword hash user password
func (r *Registration) hashPassword() error {
	hash, err := bcrypt.GenerateFromPassword([]byte(r.password), 10)
	if err != nil {
		return err
	}

	r.password = string(hash)

	return nil
}

func (r *Registration) Events() []mediator.Event {
	return r.BaseAggregateRoot.Events()
}
