package v1

import (
	"context"

	"github.com/segmentio/kafka-go"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"

	"github.com/KyKyPy3/clean/internal/application/core"
	"github.com/KyKyPy3/clean/internal/domain/common"
	"github.com/KyKyPy3/clean/internal/infrastructure/queue"
	reg_event "github.com/KyKyPy3/clean/internal/modules/registration/application/event"
	"github.com/KyKyPy3/clean/pkg/logger"
)

type CommandBus interface {
	Dispatch(context.Context, core.Command) error
}

type RegistrationEvents struct {
	logger   logger.Logger
	commands CommandBus
	consumer *queue.Consumer
	tracer   trace.Tracer
}

func NewRegistrationEvents(consumer *queue.Consumer, commands CommandBus, logger logger.Logger) {
	handlers := &RegistrationEvents{
		logger:   logger,
		consumer: consumer,
		commands: commands,
		tracer:   otel.Tracer(""),
	}

	consumer.Subscribe("registration", handlers.Handle)
}

func (r *RegistrationEvents) Handle(ctx context.Context, event *kafka.Message) error {
	ctx, span := r.tracer.Start(ctx, "RegistrationEvents.Handle")
	defer span.End()

	r.logger.Debugf("Receive event from queue %v", event)

	email, err := common.NewEmail("zi81@nm.ru")
	if err != nil {
		return err
	}

	cmd := reg_event.SendEmailCommand{
		ID:    "1",
		Email: email,
	}
	err = r.commands.Dispatch(ctx, cmd)
	if err != nil {
		r.logger.Errorf("Can't execute send email command, err: %v", err)
	}

	return nil
}
