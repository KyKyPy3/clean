package postgres

import (
	"time"

	"github.com/KyKyPy3/clean/internal/user/domain/entity"
	"github.com/google/uuid"
)

// Database user representation
type DbUser struct {
	ID         uuid.UUID `db:"id"`
	Name       string    `db:"name"`
	Surname    string    `db:"surname"`
	Middlename string    `db:"middlename"`
	Email      string    `db:"email"`
	CreatedAt  time.Time `db:"created_at"`
	UpdatedAt  time.Time `db:"updated_at"`
}

// Convert database user model to domain model
func UserFromDB(dbUser DbUser) entity.User {
	return entity.User{
		ID:         dbUser.ID,
		Name:       dbUser.Name,
		Surname:    dbUser.Surname,
		Middlename: dbUser.Middlename,
		Email:      dbUser.Email,
		CreatedAt:  dbUser.CreatedAt,
		UpdatedAt:  dbUser.UpdatedAt,
	}
}

// Convert domain user model to database model
func UserToDB(user entity.User) DbUser {
	return DbUser{
		ID:         user.ID,
		Name:       user.Name,
		Surname:    user.Surname,
		Middlename: user.Middlename,
		Email:      user.Email,
	}
}
