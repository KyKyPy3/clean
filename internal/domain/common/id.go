package common

import (
	"errors"

	"github.com/google/uuid"
)

var (
	ErrEmptyUID     = errors.New("uid cannot be empty")
	ErrUIDBadFormat = errors.New("invalid uid format")
)

type UID struct {
	id uuid.UUID
}

// NewUID generates a new UID.
func NewUID() UID {
	return UID{
		id: uuid.New(),
	}
}

// NewWithSpecifiedID generates a new UID from UUID.
func NewWithSpecifiedID(id uuid.UUID) UID {
	return UID{
		id,
	}
}

// ParseUID constructs a UID from the given string.
// It returns an error if the string is not a valid UID.
func ParseUID(s string) (UID, error) {
	if len(s) == 0 {
		return UID{}, ErrEmptyUID
	}

	id, err := uuid.Parse(s)
	if err != nil {
		return UID{}, ErrUIDBadFormat
	}
	c := UID{
		id: id,
	}

	return c, nil
}

func (u UID) GetID() uuid.UUID {
	return u.id
}

// IsEmpty checks if the ID is nil.
// A ID is considered nil if its underlying value is uuid.Nil.
func (u UID) IsEmpty() bool {
	return u.id == uuid.Nil
}

// String returns the string representation of the ID.
// If the ID is nil, an empty string is returned.
func (u UID) String() string {
	if u.id == uuid.Nil {
		return ``
	}

	return u.id.String()
}
