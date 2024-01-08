package dto

import (
	"github.com/KyKyPy3/clean/internal/domain/common"
	"github.com/KyKyPy3/clean/internal/modules/user/domain/entity"
	"github.com/KyKyPy3/clean/internal/modules/user/domain/value_object"
)

type UserDTO struct {
	Name       string `json:"name"`
	Surname    string `json:"surname"`
	Middlename string `json:"middlename"`
}

type FetchUsersDTO struct {
	Limit  int64 `query:"limit" validate:"gte=0,lte=1000"`
	Offset int64 `query:"offset" validate:"gte=0,lte=1000"`
}

type CreateUserDTO struct {
	Name       string `json:"name" validate:"required"`
	Surname    string `json:"surname"`
	Middlename string `json:"middlename"`
	Email      string `json:"email" validate:"required"`
	Password   string `json:"password" validate:"required"`
}

// Convert database user model to domain model
func UserFromRequest(reqUser CreateUserDTO) (entity.User, error) {
	u := entity.User{}
	fullName, err := value_object.NewFullName(reqUser.Name, reqUser.Surname, reqUser.Middlename)
	if err != nil {
		return u, err
	}

	email, err := common.NewEmail(reqUser.Email)
	if err != nil {
		return u, err
	}

	u.SetFullName(fullName)
	u.SetEmail(email)
	u.SetID(common.NewUID())

	return u, nil
}

// Convert domain user model to response model
func UserToResponse(user entity.User) UserDTO {
	return UserDTO{
		Name:       user.GetFullName().FirstName(),
		Surname:    user.GetFullName().LastName(),
		Middlename: user.GetFullName().MiddleName(),
	}
}
