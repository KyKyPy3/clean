package ports

import (
	"context"

	"github.com/KyKyPy3/clean/internal/domain/common"
	"github.com/KyKyPy3/clean/internal/modules/registration/domain/entity"
	user_domain "github.com/KyKyPy3/clean/internal/modules/user/domain/entity"
	"github.com/KyKyPy3/clean/pkg/mediator"
)

type Mediator interface {
	Publish(ctx context.Context, events ...mediator.Event) error
}

type TrManager interface {
	Do(ctx context.Context, fn func(ctx context.Context) error) (err error)
}

type UserViewPgStorage interface {
	GetByEmail(ctx context.Context, email common.Email) (user_domain.User, error)
}

type EmailSender interface {
	Send(ctx context.Context, destination common.Email, subject, body string) error
}

type UniquenessPolicer interface {
	IsUnique(ctx context.Context, email common.Email) (bool, error)
}

type RegistrationPgStorage interface {
	Create(ctx context.Context, registration entity.Registration) error
	Update(ctx context.Context, registration entity.Registration) error
	GetByID(ctx context.Context, id common.UID) (entity.Registration, error)
}
