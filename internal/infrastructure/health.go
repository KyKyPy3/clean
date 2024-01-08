package infrastructure

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/KyKyPy3/clean/pkg/logger"
)

type HealthHandlers struct {
	logger logger.Logger
	app    *App
}

func NewHealthHandlers(mount *echo.Group, logger logger.Logger, app *App) {
	handlers := &HealthHandlers{logger: logger, app: app}

	mount.GET("/health", handlers.healthHandler)
}

func (h *HealthHandlers) healthHandler(c echo.Context) error {
	if err := h.app.pgClient.Ping(); err != nil {
		h.logger.Errorf("Health check failed: Postgres unavailable")
		return c.JSON(http.StatusServiceUnavailable, map[string]string{"status": "Postgres is not healthy"})
	}

	if _, err := h.app.redisClient.Ping(context.Background()).Result(); err != nil {
		h.logger.Errorf("Health check failed: Redis unavailable")
		return c.JSON(http.StatusServiceUnavailable, map[string]string{"status": "Redis is not healthy"})
	}

	if _, err := h.app.kafkaClient.Brokers(); err != nil {
		h.logger.Errorf("Health check failed: Kafka unavailable. Err: %w", err)
		return c.JSON(http.StatusServiceUnavailable, map[string]string{"status": "Kafka is not healthy"})
	}

	h.logger.Infof("Health check RequestID: %s", c.Response().Header().Get(echo.HeaderXRequestID))
	return c.JSON(http.StatusOK, map[string]string{"status": "OK"})
}
