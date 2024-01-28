package ports

import (
	"context"

	"github.com/KyKyPy3/clean/internal/domain/common"
	"github.com/KyKyPy3/clean/internal/modules/registration/domain/entity"
)

type UniquenessPolicer interface {
	IsUnique(ctx context.Context, email common.Email) (bool, error)
}

type RegistrationPgStorage interface {
	Create(ctx context.Context, registration entity.Registration) error
	Update(ctx context.Context, registration entity.Registration) error
	GetByID(ctx context.Context, id common.UID) (entity.Registration, error)
}
