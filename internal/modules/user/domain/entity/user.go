package entity

import (
	"fmt"
	"github.com/KyKyPy3/clean/pkg/mediator"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/KyKyPy3/clean/internal/domain/common"
	"github.com/KyKyPy3/clean/internal/domain/core"
	"github.com/KyKyPy3/clean/internal/modules/user/domain/value_object"
)

// User struct
type User struct {
	*core.BaseAggregateRoot

	id        common.UID
	fullName  value_object.FullName
	email     common.Email
	password  string
	createdAt time.Time
	updatedAt time.Time
}

// NewUser - creates a new User instance with the provided username, password, and email.
func NewUser(fullName value_object.FullName, email common.Email, password string) (User, error) {
	if fullName.IsEmpty() {
		return User{}, fmt.Errorf("user fullname is empty, err: %w", core.ErrInvalidEntity)
	}

	if email.IsEmpty() {
		return User{}, fmt.Errorf("user email is empty, err: %w", core.ErrInvalidEntity)
	}

	// TODO: добавить проверку на уникальность email через policy

	user := User{
		BaseAggregateRoot: &core.BaseAggregateRoot{},
		id:                common.NewUID(),
		fullName:          fullName,
		email:             email,
		password:          strings.TrimSpace(password),
	}

	return user, nil
}

func Hadrate(
	id common.UID,
	fullName value_object.FullName,
	email common.Email,
	password string,
	createdAt time.Time,
	updatedAt time.Time,
) User {
	user := User{
		BaseAggregateRoot: &core.BaseAggregateRoot{},
		id:                id,
		fullName:          fullName,
		email:             email,
		password:          password,
		createdAt:         createdAt,
		updatedAt:         updatedAt,
	}

	return user
}

func (u *User) ID() common.UID {
	return u.id
}

func (u *User) IsEmpty() bool {
	return *u == User{}
}

// FullName returns the fullname of the user.
func (u *User) FullName() value_object.FullName {
	return u.fullName
}

// UpdateFullName set the fullname of the user.
func (u *User) UpdateFullName(fullName value_object.FullName) {
	u.fullName = fullName
}

// Email returns the email of the user.
func (u *User) Email() common.Email {
	return u.email
}

// UpdateEmail set the email of the user.
func (u *User) UpdateEmail(email common.Email) {
	u.email = email
}

// Password returns the password of the user.
func (u *User) Password() string {
	return u.password
}

// UpdatePassword set the password of the user.
func (u *User) UpdatePassword(password string) {
	u.password = password
}

func (u *User) CreatedAt() time.Time {
	return u.createdAt
}

func (u *User) UpdatedAt() time.Time {
	return u.updatedAt
}

// ValidatePassword compare provided password with stored
func (u *User) ValidatePassword(password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(u.password), []byte(password))
	if err != nil {
		return err
	}
	return nil
}

// String returns the string representation of the user.
func (u *User) String() string {
	return fmt.Sprintf(
		"User{ID: %s, FullName: %s, Email: %s}",
		u.ID(),
		u.FullName().String(),
		u.Email().String(),
	)
}

func (u *User) Events() []mediator.Event {
	return u.BaseAggregateRoot.Events()
}
