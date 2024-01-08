package postgres

import (
	"github.com/KyKyPy3/clean/internal/domain/common"
	"github.com/KyKyPy3/clean/internal/modules/registration/domain/entity"
)

// DbRegistration Database registration representation
type DbRegistration struct {
	ID       string `db:"id"`
	Email    string `db:"email"`
	Verified bool   `db:"verified"`
}

// RegistrationFromDB Convert database registration model to domain model
func RegistrationFromDB(dbRegistration DbRegistration) (entity.Registration, error) {
	r := entity.Registration{}

	entityID, err := common.ParseUID(dbRegistration.ID)
	if err != nil {
		return r, err
	}

	email, err := common.NewEmail(dbRegistration.Email)
	if err != nil {
		return r, err
	}

	r.SetID(entityID)
	r.SetEmail(email)
	r.SetVerified(dbRegistration.Verified)

	return r, nil
}

// RegistrationToDB Convert domain registration model to database model
func RegistrationToDB(registration entity.Registration) DbRegistration {
	return DbRegistration{
		ID:       registration.GetID().String(),
		Email:    registration.GetEmail().String(),
		Verified: registration.GetVerified(),
	}
}
