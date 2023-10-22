package common

import "errors"

// ErrNotFound not found
var ErrNotFound = errors.New("not found")

// ErrInvalidEntity invalid entity
var ErrInvalidEntity = errors.New("invalid entity")

// ErrEntityExist entity already exist
var ErrEntityExist = errors.New("entity already exist")
