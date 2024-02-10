package queue

import (
	"context"
	"time"

	"github.com/segmentio/kafka-go"

	kafkaClient "github.com/KyKyPy3/clean/pkg/kafka"
	"github.com/KyKyPy3/clean/pkg/outbox"
)

type Queue interface {
	Publish(ctx context.Context, event outbox.Message) error
}

type queue struct {
	producer kafkaClient.Producer
}

func NewQueue(producer kafkaClient.Producer) Queue {
	return &queue{
		producer: producer,
	}
}

func (q *queue) Publish(ctx context.Context, event outbox.Message) error {
	err := q.producer.PublishMessage(ctx, kafka.Message{
		Topic: event.Topic,
		Value: event.Payload,
		Time:  time.Now().UTC(),
		Headers: []kafka.Header{{
			Key:   "Kind",
			Value: []byte(event.Kind),
		}},
	})
	if err != nil {
		return err
	}

	return nil
}
