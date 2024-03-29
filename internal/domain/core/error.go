package core

import "errors"

// ErrNotFound not found.
var ErrNotFound = errors.New("not found")

// ErrInvalidEntity invalid entity.
var ErrInvalidEntity = errors.New("invalid entity")

// ErrAlreadyExist already exist.
var ErrAlreadyExist = errors.New("already exist")

// ErrNoChanges no changes.
var ErrNoChanges = errors.New("entity has not changes")
