package infrastructure

import (
	"context"
	"net/http"

	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
	"github.com/segmentio/kafka-go"

	"github.com/KyKyPy3/clean/pkg/logger"
)

type HealthHandlers struct {
	logger         logger.Logger
	postgresClient *sqlx.DB
	redisClient    *redis.Client
	kafkaClient    *kafka.Conn
}

func NewHealthHandlers(
	mount *echo.Group,
	logger logger.Logger,
	postgresClient *sqlx.DB,
	redisClient *redis.Client,
	kafkaClient *kafka.Conn,
) {
	handlers := &HealthHandlers{
		logger:         logger,
		postgresClient: postgresClient,
		redisClient:    redisClient,
		kafkaClient:    kafkaClient,
	}

	mount.GET("/health", handlers.healthHandler)
}

func (h *HealthHandlers) healthHandler(c echo.Context) error {
	if err := h.postgresClient.Ping(); err != nil {
		h.logger.Errorf("Health check failed: Postgres unavailable")
		return c.JSON(http.StatusServiceUnavailable, map[string]string{"status": "Postgres is not healthy"})
	}

	if _, err := h.redisClient.Ping(context.Background()).Result(); err != nil {
		h.logger.Errorf("Health check failed: Redis unavailable")
		return c.JSON(http.StatusServiceUnavailable, map[string]string{"status": "Redis is not healthy"})
	}

	if _, err := h.kafkaClient.Brokers(); err != nil {
		h.logger.Errorf("Health check failed: Kafka unavailable. Err: %w", err)
		return c.JSON(http.StatusServiceUnavailable, map[string]string{"status": "Kafka is not healthy"})
	}

	h.logger.Infof("Health check RequestID: %s", c.Response().Header().Get(echo.HeaderXRequestID))
	return c.JSON(http.StatusOK, map[string]string{"status": "OK"})
}
