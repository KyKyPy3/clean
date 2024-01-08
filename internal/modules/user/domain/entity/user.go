package entity

import (
	"fmt"
	"strings"
	"time"

	"github.com/KyKyPy3/clean/internal/domain"
	"github.com/KyKyPy3/clean/internal/domain/common"
	"github.com/KyKyPy3/clean/internal/modules/user/domain/value_object"
	"golang.org/x/crypto/bcrypt"
)

// User struct
type User struct {
	id common.UID

	fullName  value_object.FullName
	email     common.Email
	password  string
	createdAt time.Time
	updatedAt time.Time
}

// NewUser - creates a new User instance with the provided username, password, and email.
func NewUser(fullName value_object.FullName, email common.Email, password string) (User, error) {
	if fullName.IsEmpty() {
		return User{}, fmt.Errorf("user fullname is empty, err: %w", domain.ErrInvalidEntity)
	}

	if email.IsEmpty() {
		return User{}, fmt.Errorf("user email is empty, err: %w", domain.ErrInvalidEntity)
	}

	user := User{
		fullName: fullName,
		email:    email,
		password: strings.TrimSpace(password),
	}
	user.SetID(common.NewUID())
	if err := user.HashPassword(); err != nil {
		return User{}, err
	}

	return user, nil
}

func (u *User) GetID() common.UID {
	return u.id
}

func (u *User) SetID(id common.UID) {
	u.id = id
}

func (u *User) IsEmpty() bool {
	return *u == User{}
}

// GetFullName returns the fullname of the user.
func (u *User) GetFullName() value_object.FullName {
	return u.fullName
}

// SetFullName set the fullname of the user.
func (u *User) SetFullName(fullName value_object.FullName) {
	u.fullName = fullName
}

// GetEmail returns the email of the user.
func (u *User) GetEmail() common.Email {
	return u.email
}

// SetEmail set the email of the user.
func (u *User) SetEmail(email common.Email) {
	u.email = email
}

// GetPassword returns the password of the user.
func (u *User) GetPassword() string {
	return u.password
}

// SetPassword set the password of the user.
func (u *User) SetPassword(password string) {
	u.password = password
}

func (u *User) GetCreatedAt() time.Time {
	return u.createdAt
}

func (u *User) SetCreatedAt(createdAt time.Time) {
	u.createdAt = createdAt
}

func (u *User) GetUpdatedAt() time.Time {
	return u.updatedAt
}

func (u *User) SetUpdatedAt(updatedAt time.Time) {
	u.updatedAt = updatedAt
}

// ValidatePassword compare provided password with stored
func (u *User) ValidatePassword(password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(u.password), []byte(password))
	if err != nil {
		return err
	}
	return nil
}

// HashPassword hash user password
func (u *User) HashPassword() error {
	hash, err := bcrypt.GenerateFromPassword([]byte(u.password), 10)
	if err != nil {
		return err
	}

	u.password = string(hash)

	return nil
}

// String returns the string representation of the user.
func (u *User) String() string {
	return fmt.Sprintf(
		"User{ID: %s, FullName: %s, Email: %s}",
		u.GetID(),
		u.GetFullName().String(),
		u.GetEmail().String(),
	)
}
