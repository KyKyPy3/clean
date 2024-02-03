package registration

import (
	"context"

	trmsqlx "github.com/avito-tech/go-transaction-manager/drivers/sqlx/v2"
	"github.com/avito-tech/go-transaction-manager/trm/v2/manager"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"

	"github.com/KyKyPy3/clean/internal/application/core"
	"github.com/KyKyPy3/clean/internal/infrastructure/queue"
	"github.com/KyKyPy3/clean/internal/modules/registration/application"
	"github.com/KyKyPy3/clean/internal/modules/registration/application/command"
	reg_event "github.com/KyKyPy3/clean/internal/modules/registration/application/event"
	"github.com/KyKyPy3/clean/internal/modules/registration/domain/event"
	handlers "github.com/KyKyPy3/clean/internal/modules/registration/infrastructure/controller/http/v1"
	events "github.com/KyKyPy3/clean/internal/modules/registration/infrastructure/controller/queue/v1"
	"github.com/KyKyPy3/clean/internal/modules/registration/infrastructure/gateway/email"
	registrationPg "github.com/KyKyPy3/clean/internal/modules/registration/infrastructure/gateway/postgres"
	user_event "github.com/KyKyPy3/clean/internal/modules/user/application/event"
	userPg "github.com/KyKyPy3/clean/internal/modules/user/infrastructure/gateway/postgres"
	"github.com/KyKyPy3/clean/pkg/logger"
	"github.com/KyKyPy3/clean/pkg/mediator"
	"github.com/KyKyPy3/clean/pkg/outbox"
)

const (
	queueTopic = "registration"
)

func InitUserHandlers(
	pgClient *sqlx.DB,
	mountPoint *echo.Group,
	pubsub *mediator.Mediator,
	consumer *queue.Consumer,
	trManager *manager.Manager,
	emailGateway *email.Client,
	outboxManager outbox.Manager,
	logger logger.Logger,
) {
	userPgStorage := userPg.NewUserPgStorage(pgClient, trmsqlx.DefaultCtxGetter, logger)
	regPgStorage := registrationPg.NewRegistrationPgStorage(pgClient, trmsqlx.DefaultCtxGetter, logger)
	regUniqPolicy := application.NewUniquenessPolicy(userPgStorage, logger)
	regCmdBus := core.NewCommandBus()
	regCmdBus.Register(
		command.CreateRegistrationKind,
		command.NewCreateRegistration(regPgStorage, regUniqPolicy, trManager, pubsub, logger),
	)
	regCmdBus.Register(
		command.ConfirmRegistrationKind,
		command.NewConfirmRegistration(regPgStorage, pubsub, trManager, logger),
	)
	regCmdBus.Register(
		reg_event.SendEmailKind,
		reg_event.NewSendEmail(logger, emailGateway),
	)

	pubsub.Subscribe(event.RegistrationCreated, func(ctx context.Context, e mediator.Event) error {
		logger.Debugf("Receive domain event %v", e)

		err := outboxManager.Publish(ctx, queueTopic, e)
		if err != nil {
			return err
		}

		return nil
	})

	pubsub.Subscribe(event.RegistrationVerified, user_event.NewRegistrationVerified(logger, userPgStorage).Handle)

	handlers.NewRegistrationHandlers(mountPoint, regCmdBus, logger)
	events.NewRegistrationEvents(consumer, regCmdBus, logger)
}
