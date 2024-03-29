package event

import "github.com/KyKyPy3/clean/internal/domain/common"

const RegistrationVerified = "RegistrationVerified"

type RegistrationVerifiedEvent struct {
	ID       string
	Email    common.Email
	Password string
}

func (e RegistrationVerifiedEvent) Kind() string {
	return RegistrationVerified
}
