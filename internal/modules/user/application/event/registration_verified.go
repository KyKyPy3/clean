package event

import (
	"context"
	"fmt"
	"github.com/KyKyPy3/clean/internal/modules/user/domain/value_object"
	"strings"

	"github.com/KyKyPy3/clean/internal/modules/registration/domain/event"
	"github.com/KyKyPy3/clean/internal/modules/user/application/ports"
	"github.com/KyKyPy3/clean/internal/modules/user/domain/entity"
	"github.com/KyKyPy3/clean/pkg/logger"
	"github.com/KyKyPy3/clean/pkg/mediator"
)

type RegistrationVerified struct {
	storage ports.UserPgStorage
	policy  ports.UniquenessPolicer
	logger  logger.Logger
}

func NewRegistrationVerified(logger logger.Logger, storage ports.UserPgStorage, policy ports.UniquenessPolicer) *RegistrationVerified {
	return &RegistrationVerified{
		storage: storage,
		policy:  policy,
		logger:  logger,
	}
}

func (r *RegistrationVerified) Handle(ctx context.Context, e mediator.Event) error {
	switch t := e.(type) {
	case event.RegistrationVerifiedEvent:
		return r.handleEmailVerified(ctx, t)
	default:
		return fmt.Errorf("Unknown type of event %T", e)
	}
}

func (r RegistrationVerified) handleEmailVerified(ctx context.Context, e event.RegistrationVerifiedEvent) error {
	r.logger.Debugf("Get email verified event")

	fullName, err := value_object.NewFullName(strings.Split(e.Email.String(), "@")[0], "", "")
	if err != nil {
		return err
	}

	user, err := entity.NewUser(fullName, e.Email, e.Password, r.policy)
	if err != nil {
		return err
	}

	r.logger.Debugf("Create user %+v", user)

	err = r.storage.Create(ctx, user)
	if err != nil {
		return err
	}

	return nil
}
