package v1

import (
	"context"
	"encoding/json"

	"github.com/segmentio/kafka-go"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"

	"github.com/KyKyPy3/clean/internal/application/core"
	"github.com/KyKyPy3/clean/internal/infrastructure/queue"
	reg_event "github.com/KyKyPy3/clean/internal/modules/registration/application/event"
	"github.com/KyKyPy3/clean/pkg/logger"
)

type CommandBus interface {
	Dispatch(context.Context, core.Command) (any, error)
}

type RegistrationEvent struct {
	ID    string `json:"id"`
	Email string `json:"email"`
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

	regEvent := RegistrationEvent{}
	err := json.Unmarshal(event.Value, &regEvent)
	if err != nil {
		return err
	}

	r.logger.Debugf("Receive event from queue %+v", regEvent)

	cmd := reg_event.SendEmailCommand{
		ID:    regEvent.ID,
		Email: regEvent.Email,
	}
	_, err = r.commands.Dispatch(ctx, cmd)
	if err != nil {
		r.logger.Errorf("Can't execute send email command, err: %v", err)
	}

	return nil
}
