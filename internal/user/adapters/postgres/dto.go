package postgres

import (
	"time"

	"github.com/KyKyPy3/clean/internal/user/domain/entity"
	"github.com/google/uuid"
)

// DbUser Database user representation
type DbUser struct {
	ID         uuid.UUID `db:"id"`
	Name       string    `db:"name"`
	Surname    string    `db:"surname"`
	Middlename string    `db:"middlename"`
	Email      string    `db:"email"`
	CreatedAt  time.Time `db:"created_at"`
	UpdatedAt  time.Time `db:"updated_at"`
}

// UserFromDB Convert database user model to domain model
func UserFromDB(dbUser DbUser) entity.User {
	u := entity.User{}
	u.SetID(dbUser.ID)
	u.SetFirstName(dbUser.Name)
	u.SetLastName(dbUser.Surname)
	u.SetMiddleName(dbUser.Middlename)
	u.SetEmail(dbUser.Email)
	u.SetCreatedAt(dbUser.CreatedAt)
	u.SetUpdatedAt(dbUser.CreatedAt)

	return u
}

// UserToDB Convert domain user model to database model
func UserToDB(user entity.User) DbUser {
	return DbUser{
		ID:         user.ID(),
		Name:       user.FirstName(),
		Surname:    user.LastName(),
		Middlename: user.MiddleName(),
		Email:      user.Email(),
	}
}
