package postgres

import (
	"time"

	"github.com/KyKyPy3/clean/internal/domain/common"
	"github.com/KyKyPy3/clean/internal/modules/user/domain/entity"
	"github.com/KyKyPy3/clean/internal/modules/user/domain/value_object"
)

// DbUser Database user representation
type DbUser struct {
	ID         string    `db:"id"`
	Name       string    `db:"name"`
	Surname    string    `db:"surname"`
	Middlename string    `db:"middlename"`
	Email      string    `db:"email"`
	CreatedAt  time.Time `db:"created_at"`
	UpdatedAt  time.Time `db:"updated_at"`
}

// UserFromDB Convert database user model to domain model
func UserFromDB(dbUser DbUser) (entity.User, error) {
	u := entity.User{}

	entityID, err := common.ParseUID(dbUser.ID)
	if err != nil {
		return u, err
	}

	email, err := common.NewEmail(dbUser.Email)
	if err != nil {
		return u, err
	}

	fullName, err := value_object.NewFullName(dbUser.Name, dbUser.Surname, dbUser.Middlename)
	if err != nil {
		return u, err
	}

	u.SetID(entityID)
	u.SetFullName(fullName)
	u.SetEmail(email)
	u.SetCreatedAt(dbUser.CreatedAt)
	u.SetUpdatedAt(dbUser.CreatedAt)

	return u, nil
}

// UserToDB Convert domain user model to database model
func UserToDB(user entity.User) DbUser {
	return DbUser{
		ID:         user.GetID().String(),
		Name:       user.GetFullName().FirstName(),
		Surname:    user.GetFullName().LastName(),
		Middlename: user.GetFullName().MiddleName(),
		Email:      user.GetEmail().String(),
	}
}
