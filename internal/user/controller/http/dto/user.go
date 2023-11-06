package dto

import "github.com/KyKyPy3/clean/internal/user/domain/entity"

type UserDTO struct {
	Name       string `json:"name"`
	Surname    string `json:"surname"`
	Middlename string `json:"middlename"`
}

type ValidationError struct {
	Field  string      `json:"field"`
	Value  interface{} `json:"value"`
	Reason string      `json:"reason"`
}

type ResponseDTO struct {
	Status  int                `json:"status"`
	Message string             `json:"message"`
	Data    interface{}        `json:"data,omitempty"`
	Errors  []*ValidationError `json:"errors,omitempty"`
	Error   string             `json:"error,omitempty"`
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
func UserFromRequest(reqUser CreateUserDTO) entity.User {
	return entity.User{
		Name:       reqUser.Name,
		Surname:    reqUser.Surname,
		Middlename: reqUser.Middlename,
		Email:      reqUser.Email,
	}
}

// Convert domain user model to response model
func UserToResponse(user entity.User) UserDTO {
	return UserDTO{
		Name:       user.Name,
		Surname:    user.Surname,
		Middlename: user.Middlename,
	}
}
