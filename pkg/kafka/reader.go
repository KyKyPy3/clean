package kafka

import (
	"time"

	"github.com/segmentio/kafka-go"
)

const (
	maxAttempts            = 10
	maxWaitTime            = 3 * time.Second
	heartbeatInterval      = 1 * time.Second
	minBytes               = 10e3 // 10KB
	maxBytes               = 10e6 // 10MB
	queueCapacity          = 100
	commitInterval         = 0
	partitionWatchInterval = 500 * time.Millisecond
	dialTimeout            = 3 * time.Minute
	maxReadBackoff         = 300 * time.Millisecond
)

// NewKafkaReader create new configured kafka reader.
func NewKafkaReader(kafkaURL []string, topic, groupID string, errLogger kafka.Logger) *kafka.Reader {
	return kafka.NewReader(kafka.ReaderConfig{
		Brokers:                kafkaURL,
		GroupID:                groupID,
		Topic:                  topic,
		MinBytes:               minBytes,
		MaxBytes:               maxBytes,
		QueueCapacity:          queueCapacity,
		HeartbeatInterval:      heartbeatInterval,
		CommitInterval:         commitInterval,
		PartitionWatchInterval: partitionWatchInterval,
		ErrorLogger:            errLogger,
		MaxAttempts:            maxAttempts,
		MaxWait:                maxWaitTime,
		Dialer:                 &kafka.Dialer{Timeout: dialTimeout},
		ReadBackoffMax:         maxReadBackoff,
	})
}
