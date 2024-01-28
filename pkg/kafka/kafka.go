package kafka

import (
	"context"
	"github.com/segmentio/kafka-go"
)

type Config struct {
	Brokers []string
	GroupID string
}

func New(ctx context.Context, config Config) (*kafka.Conn, error) {
	return kafka.DialContext(ctx, "tcp", config.Brokers[0])
}
