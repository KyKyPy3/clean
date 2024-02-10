package event

import (
	"github.com/KyKyPy3/clean/internal/domain/common"
	"github.com/KyKyPy3/clean/internal/modules/user/domain/vo"
)

const UserCreated = "UserCreated"

type UserCreatedEvent struct {
	ID       string
	FullName vo.FullName
	Email    common.Email
}

func (e UserCreatedEvent) Kind() string {
	return UserCreated
}
