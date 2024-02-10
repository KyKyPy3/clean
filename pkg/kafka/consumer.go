package kafka

import (
	"context"
	"time"

	"github.com/segmentio/kafka-go"
	"golang.org/x/sync/errgroup"

	"github.com/KyKyPy3/clean/pkg/logger"
)

const maxWait = 1 * time.Second

type Consumer interface {
	ConsumeTopic(ctx context.Context, topics []string, poolSize int, worker Worker) error
}

// Worker kafka consumer worker fetch and process messages from reader.
type Worker func(ctx context.Context, r *kafka.Reader, workerID int) error

type consumer struct {
	Brokers []string
	GroupID string
	log     logger.Logger
}

// NewConsumer kafka consumer constructor.
func NewConsumer(brokers []string, groupID string, log logger.Logger) Consumer {
	return &consumer{Brokers: brokers, GroupID: groupID, log: log}
}

// GetNewKafkaReader create new kafka reader.
func (c *consumer) GetNewKafkaReader(kafkaURL []string, groupTopics []string, groupID string) *kafka.Reader {
	return kafka.NewReader(kafka.ReaderConfig{
		Brokers:                kafkaURL,
		GroupID:                groupID,
		GroupTopics:            groupTopics,
		MinBytes:               minBytes,
		MaxBytes:               maxBytes,
		QueueCapacity:          queueCapacity,
		HeartbeatInterval:      heartbeatInterval,
		CommitInterval:         commitInterval,
		PartitionWatchInterval: partitionWatchInterval,
		MaxAttempts:            maxAttempts,
		MaxWait:                maxWait,
		Dialer:                 &kafka.Dialer{Timeout: dialTimeout},
	})
}

// ConsumeTopic start consumer group with given worker and pool size.
func (c *consumer) ConsumeTopic(ctx context.Context, topics []string, poolSize int, worker Worker) error {
	r := c.GetNewKafkaReader(c.Brokers, topics, c.GroupID)

	defer func() {
		if err := r.Close(); err != nil {
			c.log.Warnf("consumer.r.Close: %v", err)
		}
	}()

	c.log.Infof("(Starting consumer groupID): GroupID %s, topic: %+v, poolSize: %v", c.GroupID, topics, poolSize)

	g, ctx := errgroup.WithContext(ctx)

	for i := 0; i <= poolSize; i++ {
		g.Go(c.runWorker(ctx, worker, r, i))
	}

	return g.Wait()
}

func (c *consumer) runWorker(ctx context.Context, worker Worker, r *kafka.Reader, i int) func() error {
	return func() error {
		return worker(ctx, r, i)
	}
}
