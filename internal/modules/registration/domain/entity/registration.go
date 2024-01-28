package entity

import (
	"context"
	"fmt"
	"github.com/KyKyPy3/clean/pkg/mediator"

	"github.com/KyKyPy3/clean/internal/domain/common"
	"github.com/KyKyPy3/clean/internal/domain/core"
	"github.com/KyKyPy3/clean/internal/modules/registration/domain"
	"github.com/KyKyPy3/clean/internal/modules/registration/domain/event"
)

// Registration struct
type Registration struct {
	*core.BaseAggregateRoot

	id       common.UID
	email    common.Email
	verified bool
}

// NewRegistration - create and validate registration
func NewRegistration(ctx context.Context, email common.Email, uniqPolicy domain.UniqueEmailPolicy) (Registration, error) {
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

	agregateRoot := core.BaseAggregateRoot{}
	r := Registration{
		BaseAggregateRoot: &agregateRoot,
		id:                common.NewUID(),
		email:             email,
		verified:          false,
	}

	r.BaseAggregateRoot.AddEvent(event.RegistrationCreatedEvent{ID: r.ID().String(), Email: email})

	return r, nil
}

func Hydrate(id common.UID, email common.Email, verified bool) (Registration, error) {
	if email.IsEmpty() {
		return Registration{}, fmt.Errorf("registration email is empty, err: %w", core.ErrInvalidEntity)
	}

	agregateRoot := core.BaseAggregateRoot{}
	reg := Registration{
		BaseAggregateRoot: &agregateRoot,
		id:                id,
		email:             email,
		verified:          verified,
	}

	return reg, nil
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

func (r *Registration) Verified() bool {
	return r.verified
}

func (r *Registration) Verify() error {
	if r.verified {
		return core.ErrNoChanges
	}

	r.verified = true
	r.BaseAggregateRoot.AddEvent(event.RegistrationVerifiedEvent{ID: r.ID().String(), Email: r.Email()})

	return nil
}

func (r *Registration) Events() []mediator.Event {
	return r.BaseAggregateRoot.Events()
}
