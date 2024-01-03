package valueobject

import (
	"errors"
	"regexp"
)

var emailRe = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

var (
	ErrEmptyEmail = errors.New("email cannot be empty")
	ErrBadFormat = errors.New("invalid email format")
)

type Email struct {
	email string
}

func NewEmail(email string) (Email, error) {
	if email == "" {
		return Email{}, ErrEmptyEmail	
	}
	
	if !emailRe.MatchString(email) {
		return Email{}, ErrBadFormat
	}

	return Email{
		email: email,
	}, nil
}