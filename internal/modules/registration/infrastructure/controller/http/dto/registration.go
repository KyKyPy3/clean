package dto

type CreateRegistrationDTO struct {
	Email string `json:"email" validate:"required"`
}
