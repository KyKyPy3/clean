package postgres

import (
	"time"

	"github.com/KyKyPy3/clean/internal/domain/common"
	"github.com/KyKyPy3/clean/internal/modules/user/domain/entity"
	"github.com/KyKyPy3/clean/internal/modules/user/domain/vo"
)

// DBUser Database user representation.
type DBUser struct {
	ID         string    `db:"id"`
	Name       string    `db:"name"`
	Surname    string    `db:"surname"`
	Middlename string    `db:"middlename"`
	Email      string    `db:"email"`
	Password   string    `db:"password"`
	CreatedAt  time.Time `db:"created_at"`
	UpdatedAt  time.Time `db:"updated_at"`
}

// UserFromDB Convert database user model to domain model.
func UserFromDB(dbUser DBUser) (entity.User, error) {
	entityID, err := common.ParseUID(dbUser.ID)
	if err != nil {
		return entity.User{}, err
	}

	email, err := common.NewEmail(dbUser.Email)
	if err != nil {
		return entity.User{}, err
	}

	fullName, err := vo.NewFullName(dbUser.Name, dbUser.Surname, dbUser.Middlename)
	if err != nil {
		return entity.User{}, err
	}

	user := entity.Hydrate(entityID, fullName, email, dbUser.Password, dbUser.CreatedAt, dbUser.UpdatedAt)

	return user, nil
}

// UserToDB Convert domain user model to database model.
func UserToDB(user entity.User) DBUser {
	return DBUser{
		ID:         user.ID().String(),
		Name:       user.FullName().FirstName(),
		Surname:    user.FullName().LastName(),
		Middlename: user.FullName().MiddleName(),
		Email:      user.Email().String(),
		Password:   user.Password(),
	}
}
