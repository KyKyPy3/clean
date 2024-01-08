package event

import "github.com/KyKyPy3/clean/internal/domain/common"

type RegistrationCreatedEvent struct {
	ID    string
	Email common.Email
}
