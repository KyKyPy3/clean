package common

import (
	"regexp"
	"strings"
)

import (
	"errors"
)

var emailRe = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

var (
	ErrEmptyEmail = errors.New("email cannot be empty")
	ErrBadFormat  = errors.New("invalid email format")
)

type Email struct {
	email string
}

func NewEmail(email string) (Email, error) {
	filteredEmail := strings.ToLower(strings.TrimSpace(email))
	if filteredEmail == "" {
		return Email{}, ErrEmptyEmail
	}

	e := Email{
		email: filteredEmail,
	}

	if err := e.validate(); err != nil {
		return Email{}, err
	}

	return e, nil
}

func MustNewEmail(email string) Email {
	e, err := NewEmail(email)
	if err != nil {
		panic(err)
	}

	return e
}

func (e Email) validate() error {
	if !emailRe.MatchString(e.email) {
		return ErrBadFormat
	}
	return nil
}

func (e Email) IsEmpty() bool {
	return e == Email{}
}

func (e Email) String() string {
	return e.email
}
