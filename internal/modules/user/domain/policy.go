package domain

import (
	"github.com/KyKyPy3/clean/internal/domain/common"
)

type UniqueEmailPolicy interface {
	IsUnique(common.Email) (bool, error)
}
