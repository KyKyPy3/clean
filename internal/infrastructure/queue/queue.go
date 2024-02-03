package queue

import (
	"context"
	"sync"

	"github.com/segmentio/kafka-go"

	"github.com/KyKyPy3/clean/internal/infrastructure/config"
	kafkaClient "github.com/KyKyPy3/clean/pkg/kafka"
	"github.com/KyKyPy3/clean/pkg/latch"
	"github.com/KyKyPy3/clean/pkg/logger"
)

type EventHandler func(ctx context.Context, event *kafka.Message) error

type Consumer struct {
	logger   logger.Logger
	cfg      *config.Config
	handlers map[string]EventHandler
	lock     *latch.CountDownLatch
	mutex    sync.Mutex
}

func NewConsumer(cfg *config.Config, lock *latch.CountDownLatch, logger logger.Logger) *Consumer {
	return &Consumer{
		logger:   logger,
		cfg:      cfg,
		lock:     lock,
		handlers: make(map[string]EventHandler),
	}
}

func (c *Consumer) Start(ctx context.Context) error {
	kafkaConsumer := kafkaClient.NewConsumer(c.cfg.Kafka.Brokers, c.cfg.Kafka.GroupID, c.logger)

	c.lock.Add(1)
	go func() {
		defer c.lock.Done()

		err := kafkaConsumer.ConsumeTopic(
			ctx,
			[]string{"registration"},
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
						c.logger.Warnf("(kafkaConsumer ConsumeTopic) workerID: %d, err: %v", workerID, err)
						continue
					}

					c.logger.Debugf("Message %#v", msg)

					handler, ok := c.handlers[msg.Topic]
					if !ok {
						c.logger.Warnf("(kafkaConsumer ConsumeTopic) event handler for topic %s not found", msg.Topic)
					}
					err = handler(ctx, &msg)
					if err != nil {
						c.logger.Errorf("(kafkaConsumer ConsumeTopic) can't process event, err: %v", err)
					}

					if err := r.CommitMessages(ctx, msg); err != nil {
						c.logger.Errorf("failed to commit messages: %v", err)
					}
				}
			},
		)
		if err != nil {
			c.logger.Errorf("(kafkaConsumer ConsumeTopic) err: %v", err)
			return
		}
	}()

	return nil
}

func (c *Consumer) Subscribe(eventType string, handler EventHandler) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.handlers[eventType] = handler
}

func (_ *Consumer) Stop() error {
	return nil
}
