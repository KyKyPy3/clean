package dto

type UserDTO struct {
	Name string `json:"name"`
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
	Limit int64 `query:"limit" validate:"gte=0,lte=1000"`
}

type CreateUserDTO struct {
	Name       string `json:"name" validate:"required"`
	Surname    string `json:"surname"`
	Middlename string `json:"middlename"`
	Email      string `json:"email" validate:"required"`
	Password   string `json:"password" validate:"required"`
}
