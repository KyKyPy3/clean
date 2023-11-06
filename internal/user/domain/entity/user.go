package entity

import (
	"strings"
	"time"

	"github.com/KyKyPy3/clean/internal/common"
	"golang.org/x/crypto/bcrypt"
)

// User struct
type User struct {
	ID         common.ID
	Name       string
	Surname    string
	Middlename string
	Email      string
	Password   string
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

// Create, unify and validate user
func NewUser(name, surname, middlename, email, password string) (User, error) {
	u := User{
		Name:       strings.TrimSpace(name),
		Surname:    strings.TrimSpace(surname),
		Middlename: strings.TrimSpace(middlename),
		Email:      strings.ToLower(strings.TrimSpace(email)),
		Password:   strings.TrimSpace(password),
	}
	password, err := hashPassword(u.Password)
	if err != nil {
		return User{}, err
	}
	u.Password = password
	err = u.Validate()
	if err != nil {
		return User{}, err
	}

	return u, nil
}

// Compare provided password with stored
func (u *User) ValidatePassword(password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	if err != nil {
		return err
	}
	return nil
}

// Validate user
func (u *User) Validate() error {
	if u.Email == "" {
		return common.ErrInvalidEntity
	}
	return nil
}

// Hash given password
func hashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return "", err
	}

	return string(hash), nil
}
