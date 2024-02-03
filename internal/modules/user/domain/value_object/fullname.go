package value_object

import (
	"errors"
	"fmt"
	"strings"
)

var (
	ErrEmptyFirstName = errors.New("firstName cannot be empty")
)

// FullName is a value object representing first, last and middle names
type FullName struct {
	firstName  string
	lastName   string
	middleName string
}

func NewFullName(firstName, lastName, middleName string) (FullName, error) {
	fullName := FullName{
		firstName:  strings.TrimSpace(firstName),
		lastName:   strings.TrimSpace(lastName),
		middleName: strings.TrimSpace(middleName),
	}

	if fullName.firstName == "" {
		return FullName{}, ErrEmptyFirstName
	}

	return fullName, nil
}

func MustNewFullName(firstName, lastName, middleName string) FullName {
	fullName, err := NewFullName(firstName, lastName, middleName)
	if err != nil {
		panic(err)
	}

	return fullName
}

func (f FullName) IsEmpty() bool {
	return f == FullName{}
}

func (f FullName) FirstName() string {
	return f.firstName
}

func (f FullName) LastName() string {
	return f.lastName
}

func (f FullName) MiddleName() string {
	return f.middleName
}

func (f FullName) String() string {
	return fmt.Sprintf("%s %s %s", f.lastName, f.firstName, f.middleName)
}
