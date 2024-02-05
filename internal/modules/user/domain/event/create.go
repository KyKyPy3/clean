package event

import (
	"github.com/KyKyPy3/clean/internal/domain/common"
	"github.com/KyKyPy3/clean/internal/modules/user/domain/value_object"
)

const UserCreated = "UserCreated"

type UserCreatedEvent struct {
	ID       string
	FullName value_object.FullName
	Email    common.Email
}

func (e UserCreatedEvent) Kind() string {
	return UserCreated
}
