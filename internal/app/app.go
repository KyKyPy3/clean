package app

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/KyKyPy3/clean/config"
	"github.com/KyKyPy3/clean/internal/user/adapters/postgres"
	redis2 "github.com/KyKyPy3/clean/internal/user/adapters/redis"
	v1 "github.com/KyKyPy3/clean/internal/user/controller/http/v1"
	"github.com/KyKyPy3/clean/internal/user/domain/service"
	"github.com/KyKyPy3/clean/internal/user/usecase"
	"github.com/KyKyPy3/clean/pkg/logger"
	"github.com/go-playground/validator"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/redis/go-redis/v9"
	"github.com/segmentio/kafka-go"
	"go.opentelemetry.io/contrib/instrumentation/github.com/labstack/echo/otelecho"
	"go.opentelemetry.io/otel"
	otelmetric "go.opentelemetry.io/otel/metric"
)

const (
	shutdownTimeout = 30 * time.Second
	maxHeaderBytes  = 1 << 20
)

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	if err := cv.validator.Struct(i); err != nil {
		return err
	}
	return nil
}

type App struct {
	cfg         *config.Config
	echo        *echo.Echo
	logger      logger.Logger
	pgClient    *sqlx.DB
	redisClient *redis.Client
	kafkaClient *kafka.Conn
}

func New(
	cfg *config.Config,
	logger logger.Logger,
	pgClient *sqlx.DB,
	redisClient *redis.Client,
	kafkaClient *kafka.Conn,
) *App {
	return &App{
		cfg:         cfg,
		logger:      logger,
		pgClient:    pgClient,
		echo:        echo.New(),
		redisClient: redisClient,
		kafkaClient: kafkaClient,
	}
}

func (a *App) Run(ctx context.Context) error {
	// Collect metrics
	meter := otel.GetMeterProvider().Meter("")
	testMetric, _ := meter.Int64Counter(
		"test_metric",
		otelmetric.WithDescription("Number of jobs received by the compute node"),
	)
	testMetric.Add(ctx, 1)

	a.echo.Validator = &CustomValidator{validator: validator.New()}

	a.echo.Use(middleware.Logger())
	a.echo.Use(middleware.Recover())
	// a.echo.Use(middleware.CSRF())
	a.echo.Use(middleware.CORS())
	a.echo.Use(otelecho.Middleware("clean"))

	a.echo.GET("/healthz", func(c echo.Context) error {
		a.logger.Info("Health check triggered")
		return c.JSON(http.StatusOK, map[string]string{"status": "OK"})
	})

	mountPoint := a.echo.Group("/api/v1")
	a.connectHandlers(mountPoint)

	serverCtx, cancel := context.WithCancel(context.Background())

	// Listen for syscall signals for process to interrupt/quit
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		<-sig

		// Shutdown signal with grace period of 30 seconds
		shutdownCtx, _ := context.WithTimeout(serverCtx, shutdownTimeout) //nolint:all
		go func() {
			<-shutdownCtx.Done()
			if shutdownCtx.Err() == context.DeadlineExceeded {
				a.logger.Fatal("graceful shutdown timed out.. forcing exit.")
			}
		}()

		// Trigger graceful shutdown
		err := a.echo.Shutdown(shutdownCtx)
		if err != nil {
			a.logger.Fatal(err)
		}
		cancel()
	}()

	// Run the server
	server := &http.Server{
		Addr:           fmt.Sprintf("%s:%s", a.cfg.Server.Host, a.cfg.Server.Port),
		ReadTimeout:    time.Second * a.cfg.Server.ReadTimeout,
		WriteTimeout:   time.Second * a.cfg.Server.WriteTimeout,
		MaxHeaderBytes: maxHeaderBytes,
	}
	err := a.echo.StartServer(server)
	if err != nil && err != http.ErrServerClosed {
		return err
	}

	// Wait for server context to be stopped
	<-serverCtx.Done()

	return nil
}

func (a *App) connectHandlers(mountPoint *echo.Group) {
	// Health endpoint
	health := mountPoint.Group("/health")
	health.GET("", func(c echo.Context) error {
		if err := a.pgClient.Ping(); err != nil {
			a.logger.Errorf("Health check failed: Postgres unavailable")
			return c.JSON(http.StatusServiceUnavailable, map[string]string{"status": "Postgres is not healthy"})
		}

		if _, err := a.redisClient.Ping(context.Background()).Result(); err != nil {
			a.logger.Errorf("Health check failed: Redis unavailable")
			return c.JSON(http.StatusServiceUnavailable, map[string]string{"status": "Redis is not healthy"})
		}

		a.logger.Infof("Health check RequestID: %s", c.Response().Header().Get(echo.HeaderXRequestID))
		return c.JSON(http.StatusOK, map[string]string{"status": "OK"})
	})

	// Init user layers
	userPgStorage := postgres.NewUserPgStorage(a.pgClient, a.logger)
	userRedisStorage := redis2.NewUserRedisStorage(a.redisClient, a.logger)
	userService := service.NewUserService(userPgStorage, userRedisStorage, a.logger)
	userUsecase := usecase.NewUserUsecase(userService, a.logger)

	v1.NewUserHandlers(mountPoint, userUsecase, a.logger)
}
