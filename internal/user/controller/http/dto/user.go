package dto

type CreateUserDTO struct {
	Name       string `json:"name"`
	Surname    string `json:"surname"`
	Middlename string `json:"middlename"`
	Email      string `json:"email"`
	Password   string `json:"password"`
}
