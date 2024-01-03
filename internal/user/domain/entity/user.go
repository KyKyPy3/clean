package entity

import (
	"strings"
	"time"

	"github.com/KyKyPy3/clean/internal/common"
	"golang.org/x/crypto/bcrypt"
)

var _ common.Entity = (*User)(nil)

// User struct
type User struct {
	common.BaseEntity
	id         common.ID
	firstName  string
	lastName   string
	middleName string
	email      string
	password   string
	createdAt  time.Time
	updatedAt  time.Time
}

// Create, unify and validate user
func NewUser(firstName, lastName, middleName, email, password string) (User, error) {
	u := User{
		firstName:  strings.TrimSpace(firstName),
		lastName:   strings.TrimSpace(lastName),
		middleName: strings.TrimSpace(middleName),
		email:      strings.ToLower(strings.TrimSpace(email)),
		password:   strings.TrimSpace(password),
	}
	if err := u.HashPassword(); err != nil {
		return User{}, err
	}

	err := u.Validate()
	if err != nil {
		return User{}, err
	}

	return u, nil
}

func (u *User) FirstName() string {
	return u.firstName
}

func (u *User) SetFirstName(firstName string) {
	u.firstName = strings.TrimSpace(firstName)
}

func (u *User) LastName() string {
	return u.lastName
}

func (u *User) SetLastName(lastName string) {
	u.lastName = strings.TrimSpace(lastName)
}

func (u *User) MiddleName() string {
	return u.middleName
}

func (u *User) SetMiddleName(middleName string) {
	u.middleName = strings.TrimSpace(middleName)
}

func (u *User) Email() string {
	return u.email
}

func (u *User) SetEmail(email string) {
	u.email = strings.ToLower(strings.TrimSpace(email))
}

func (u *User) Password() string {
	return u.password
}

func (u *User) SetPassword(password string) {
	u.password = password
}

func (u *User) CreatedAt() time.Time {
	return u.createdAt
}

func (u *User) SetCreatedAt(createdAt time.Time) {
	u.createdAt = createdAt
}

func (u *User) UpdatedAt() time.Time {
	return u.updatedAt
}

func (u *User) SetUpdatedAt(updatedAt time.Time) {
	u.updatedAt = updatedAt
}

// Compare provided password with stored
func (u *User) ValidatePassword(password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(u.password), []byte(password))
	if err != nil {
		return err
	}
	return nil
}

func (u *User) HashPassword() error {
	hash, err := bcrypt.GenerateFromPassword([]byte(u.password), 10)
	if err != nil {
		return err
	}

	u.password = string(hash)

	return nil
}

// Validate user
func (u *User) Validate() error {
	if u.email == "" {
		return common.ErrInvalidEntity
	}
	return nil
}
