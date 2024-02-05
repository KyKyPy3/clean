package event

import "github.com/KyKyPy3/clean/internal/domain/common"

const RegistrationCreated = "RegistrationCreated"

type RegistrationCreatedEvent struct {
	ID    string
	Email common.Email
}

func (e RegistrationCreatedEvent) Kind() string {
	return RegistrationCreated
}
