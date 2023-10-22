package redis

import (
	"context"
	"fmt"
	"github.com/KyKyPy3/clean/config"
	"github.com/redis/go-redis/v9"
	"time"
)

const (
	maxRetries      = 5
	minRetryBackoff = 300 * time.Millisecond
	maxRetryBackoff = 500 * time.Millisecond
	dialTimeout     = 5 * time.Second
	readTimeout     = 5 * time.Second
	writeTimeout    = 3 * time.Second
)

func New(ctx context.Context, cfg *config.RedisConfig) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:            fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
		MinIdleConns:    cfg.MinIdleConn,
		PoolSize:        cfg.PoolSize,
		PoolTimeout:     cfg.PoolTimeout,
		Password:        cfg.Password,
		DB:              cfg.DB,
		MaxRetries:      maxRetries,
		MinRetryBackoff: minRetryBackoff,
		MaxRetryBackoff: maxRetryBackoff,
		DialTimeout:     dialTimeout,
		ReadTimeout:     readTimeout,
		WriteTimeout:    writeTimeout,
	})

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	return client, nil
}
