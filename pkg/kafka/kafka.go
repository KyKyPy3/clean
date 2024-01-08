package kafka

import (
	"context"
	"github.com/KyKyPy3/clean/internal/infrastructure/config"
	"github.com/segmentio/kafka-go"
)

func New(ctx context.Context, config *config.KafkaConfig) (*kafka.Conn, error) {
	return kafka.DialContext(ctx, "tcp", config.Brokers[0])
}
