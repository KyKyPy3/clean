package domain

import (
	"context"

	"github.com/KyKyPy3/clean/internal/domain/common"
)

type UniqueEmailPolicy interface {
	IsUnique(context.Context, common.Email) (bool, error)
}
