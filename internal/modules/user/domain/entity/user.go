package entity

import (
	"fmt"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/KyKyPy3/clean/internal/domain/common"
	"github.com/KyKyPy3/clean/internal/domain/core"
	"github.com/KyKyPy3/clean/internal/modules/user/domain"
	"github.com/KyKyPy3/clean/internal/modules/user/domain/event"
	"github.com/KyKyPy3/clean/internal/modules/user/domain/vo"
	"github.com/KyKyPy3/clean/pkg/mediator"
)

const passwordCost = 10

// User struct.
type User struct {
	*core.BaseAggregateRoot

	id        common.UID
	fullName  vo.FullName
	email     common.Email
	password  string
	createdAt time.Time
	updatedAt time.Time
}

// NewUser - creates a new User instance with the provided username, password, and email.
func NewUser(
	fullName vo.FullName,
	email common.Email,
	password string,
	uniqPolicy domain.UniqueEmailPolicy,
) (User, error) {
	if fullName.IsEmpty() {
		return User{}, fmt.Errorf("user fullname is empty, err: %w", core.ErrInvalidEntity)
	}

	if email.IsEmpty() {
		return User{}, fmt.Errorf("user email is empty, err: %w", core.ErrInvalidEntity)
	}

	ok, err := uniqPolicy.IsUnique(email)
	if err != nil {
		return User{}, fmt.Errorf("failed to check uniqueness of email on user, err: %w", err)
	}

	if !ok {
		return User{}, fmt.Errorf("user with same email already exists, err: %w", core.ErrAlreadyExist)
	}

	user := User{
		BaseAggregateRoot: &core.BaseAggregateRoot{},
		id:                common.NewUID(),
		fullName:          fullName,
		email:             email,
		password:          strings.TrimSpace(password),
	}

	user.BaseAggregateRoot.AddEvent(event.UserCreatedEvent{ID: user.ID().String(), FullName: fullName, Email: email})

	return user, nil
}

func Hydrate(
	id common.UID,
	fullName vo.FullName,
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
func (u *User) FullName() vo.FullName {
	return u.fullName
}

// UpdateFullName set the fullname of the user.
func (u *User) UpdateFullName(fullName vo.FullName) {
	u.fullName = fullName
}

// Email returns the email of the user.
func (u *User) Email() common.Email {
	return u.email
}

// UpdateEmail set the email of the user.
func (u *User) UpdateEmail(email common.Email, uniqPolicy domain.UniqueEmailPolicy) error {
	ok, err := uniqPolicy.IsUnique(email)
	if err != nil {
		return fmt.Errorf("failed to check uniqueness of email on user, err: %w", err)
	}

	if !ok {
		return fmt.Errorf("user with same email already exists, err: %w", core.ErrAlreadyExist)
	}

	u.email = email

	return nil
}

// Password returns the password of the user.
func (u *User) Password() string {
	return u.password
}

// UpdatePassword set the password of the user.
func (u *User) UpdatePassword(password string) error {
	u.password = password

	return u.hashPassword()
}

// hashPassword hash user password.
func (u *User) hashPassword() error {
	hash, err := bcrypt.GenerateFromPassword([]byte(u.password), passwordCost)
	if err != nil {
		return err
	}

	u.password = string(hash)

	return nil
}

func (u *User) CreatedAt() time.Time {
	return u.createdAt
}

func (u *User) UpdatedAt() time.Time {
	return u.updatedAt
}

// ValidatePassword compare provided password with stored.
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
