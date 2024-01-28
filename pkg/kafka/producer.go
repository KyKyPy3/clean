package kafka

import (
	"context"

	"github.com/segmentio/kafka-go"
	"go.opentelemetry.io/otel"

	"github.com/KyKyPy3/clean/pkg/logger"
)

type Producer interface {
	PublishMessage(ctx context.Context, msgs ...kafka.Message) error
	Close() error
}

type producer struct {
	log     logger.Logger
	brokers []string
	w       *kafka.Writer
}

// NewProducer create new queue producer
func NewProducer(log logger.Logger, brokers []string) Producer {
	return &producer{log: log, brokers: brokers, w: NewWriter(brokers, kafka.LoggerFunc(log.Errorf))}
}

func (p *producer) PublishMessage(ctx context.Context, msgs ...kafka.Message) error {
	ctx, span := otel.Tracer("").Start(ctx, "producer.PublishMessage")
	defer span.End()

	if err := p.w.WriteMessages(ctx, msgs...); err != nil {
		span.RecordError(err)
		return err
	}
	return nil
}

func (p *producer) Close() error {
	return p.w.Close()
}
