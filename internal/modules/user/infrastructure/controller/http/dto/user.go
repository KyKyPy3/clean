package dto

import (
	"github.com/KyKyPy3/clean/internal/modules/user/domain/entity"
)

type LoginDTO struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type UserDTO struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Surname    string `json:"surname"`
	Middlename string `json:"middlename"`
	Email      string `json:"email"`
	CreatedAt  string `json:"createdAt"`
	UpdatedAt  string `json:"updateAt"`
}

type FetchUsersDTO struct {
	Limit  int64 `query:"limit" validate:"gte=0,lte=1000"`
	Offset int64 `query:"offset" validate:"gte=0,lte=1000"`
}

type UpdateUserDTO struct {
	Name       string `json:"name" validate:"required"`
	Surname    string `json:"surname"`
	Middlename string `json:"middlename"`
	Email      string `json:"email" validate:"required"`
}

// UserToResponse - Convert domain user model to response model.
func UserToResponse(user entity.User) UserDTO {
	return UserDTO{
		ID:         user.ID().String(),
		Name:       user.FullName().FirstName(),
		Surname:    user.FullName().LastName(),
		Middlename: user.FullName().MiddleName(),
		Email:      user.Email().String(),
		CreatedAt:  user.CreatedAt().String(),
		UpdatedAt:  user.UpdatedAt().String(),
	}
}
