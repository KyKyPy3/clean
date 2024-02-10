package common

import (
	"errors"
	"regexp"
	"strings"
)

var emailRe = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$") //nolint:lll

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

func (e Email) MarshalText() ([]byte, error) {
	return []byte(e.email), nil
}

func (e Email) UnmarshalText(text []byte) error {
	e.email = string(text)
	return e.validate()
}

func (e Email) String() string {
	return e.email
}
