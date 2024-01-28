package infrastructure

import (
	"context"

	"github.com/segmentio/kafka-go"

	"github.com/KyKyPy3/clean/internal/infrastructure/config"
	kafkaClient "github.com/KyKyPy3/clean/pkg/kafka"
	"github.com/KyKyPy3/clean/pkg/latch"
	"github.com/KyKyPy3/clean/pkg/logger"
)

type Consumer struct {
	logger logger.Logger
	cfg    *config.Config
	lock   *latch.CountDownLatch
}

func NewConsumer(cfg *config.Config, lock *latch.CountDownLatch, logger logger.Logger) *Consumer {
	return &Consumer{
		logger: logger,
		cfg:    cfg,
		lock:   lock,
	}
}

func (q *Consumer) Start(ctx context.Context) error {
	kafkaConsumer := kafkaClient.NewConsumer(q.cfg.Kafka.Brokers, q.cfg.Kafka.GroupID, q.logger)

	q.lock.Add(1)
	go func() {
		defer q.lock.Done()

		err := kafkaConsumer.ConsumeTopic(
			ctx,
			[]string{"registrations"},
			5,
			func(ctx context.Context, r *kafka.Reader, workerID int) error {
				for {
					select {
					case <-ctx.Done():
						return ctx.Err()
					default:
					}

					msg, err := r.FetchMessage(ctx)
					if err != nil {
						q.logger.Warnf("(kafkaConsumer ConsumeTopic) workerID: %d, err: %v", workerID, err)
						continue
					}

					q.logger.Debugf("Message %#v", msg)
				}
			},
		)
		if err != nil {
			q.logger.Errorf("(kafkaConsumer ConsumeTopic) err: %v", err)
			return
		}
	}()

	return nil
}

func (q *Consumer) Stop(ctx context.Context) error {
	return nil
}
