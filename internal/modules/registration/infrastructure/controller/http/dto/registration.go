package dto

import (
	"github.com/KyKyPy3/clean/internal/domain/common"
	"github.com/KyKyPy3/clean/internal/modules/registration/domain/entity"
)

type RegistrationDTO struct {
	Email    string `json:"email"`
	Verified bool   `json:"verified"`
}

type CreateRegistrationDTO struct {
	Email string `json:"email" validate:"required"`
}

// Convert database registration model to domain model
func RegistrationFromRequest(reqRegistration CreateRegistrationDTO) (entity.Registration, error) {
	u := entity.Registration{}

	email, err := common.NewEmail(reqRegistration.Email)
	if err != nil {
		return u, err
	}

	u.SetEmail(email)
	u.SetID(common.NewUID())

	return u, nil
}

// Convert domain registration model to response model
func RegistrationToResponse(registration entity.Registration) RegistrationDTO {
	return RegistrationDTO{
		Email:    registration.GetEmail().String(),
		Verified: registration.GetVerified(),
	}
}
