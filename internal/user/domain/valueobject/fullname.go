package valueobject

import (
	"errors"
	"fmt"
)

var (
	ErrEmptyFirstName = errors.New("firstName cannot be empty")
	ErrEmptyLastName = errors.New("lastName cannot be empty")
)

// FullName is a value object representing first, last and middle names
type FullName struct {
	firstName string
	lastName string
	middleName string
}

func NewFullName(firstName, lastName, middleName string) (FullName, error) {
	if firstName == "" {
		return FullName{}, ErrEmptyFirstName
	}
	
	if lastName == "" {
		return FullName{}, ErrEmptyLastName
	}
	
	return FullName{
		firstName: firstName,
		lastName: lastName,
		middleName: middleName,
	}, nil
}

func (f *FullName) FirstName() string {
	return f.firstName
}

func (f *FullName) LastName() string {
	return f.lastName
}

func (f *FullName) MiddleName() string {
	return f.middleName
}

func (f *FullName) String() string {
	return fmt.Sprintf("%s %s %s", f.lastName, f.firstName, f.middleName)
}