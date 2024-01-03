package kafka

import (
	"context"
	"github.com/segmentio/kafka-go"

	"github.com/KyKyPy3/clean/config"
)

func New(ctx context.Context, config *config.KafkaConfig) (*kafka.Conn, error) {
	return kafka.DialContext(ctx, "tcp", config.Brokers[0])
}
