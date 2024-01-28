package postgres

import (
	"github.com/KyKyPy3/clean/internal/domain/common"
	"github.com/KyKyPy3/clean/internal/modules/registration/domain/entity"
)

// DBRegistration Database registration representation
type DBRegistration struct {
	ID       string `db:"id"`
	Email    string `db:"email"`
	Verified bool   `db:"verified"`
}

// RegistrationFromDB Convert database registration model to domain model
func RegistrationFromDB(dbRegistration DBRegistration) (entity.Registration, error) {
	entityID, err := common.ParseUID(dbRegistration.ID)
	if err != nil {
		return entity.Registration{}, err
	}

	email, err := common.NewEmail(dbRegistration.Email)
	if err != nil {
		return entity.Registration{}, err
	}

	r, err := entity.Hydrate(entityID, email, dbRegistration.Verified)
	if err != nil {
		return r, err
	}

	return r, nil
}

// RegistrationToDB Convert domain registration model to database model
func RegistrationToDB(registration entity.Registration) DBRegistration {
	return DBRegistration{
		ID:       registration.ID().String(),
		Email:    registration.Email().String(),
		Verified: registration.Verified(),
	}
}
