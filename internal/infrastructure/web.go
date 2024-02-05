package infrastructure

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.opentelemetry.io/contrib/instrumentation/github.com/labstack/echo/otelecho"

	"github.com/KyKyPy3/clean/internal/infrastructure/config"
	"github.com/KyKyPy3/clean/pkg/jwt"
	"github.com/KyKyPy3/clean/pkg/latch"
	"github.com/KyKyPy3/clean/pkg/logger"
)

const (
	maxHeaderBytes = 1 << 20
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

type Web struct {
	echo   *echo.Echo
	cfg    *config.Config
	logger logger.Logger
	lock   *latch.CountDownLatch
}

func NewWeb(cfg *config.Config, jwt *jwt.JWT, logger logger.Logger, lock *latch.CountDownLatch) *Web {
	return &Web{
		logger: logger,
		cfg:    cfg,
		lock:   lock,
		echo:   echo.New(),
	}
}

func (w *Web) Start() error {
	w.echo.Validator = &CustomValidator{validator: validator.New()}

	w.echo.Use(middleware.Logger())
	w.echo.Use(middleware.Recover())
	// w.echo.Use(middleware.CSRF())
	w.echo.Use(middleware.CORS())
	w.echo.Use(otelecho.Middleware("clean"))

	// Run the server
	server := &http.Server{
		Addr:           fmt.Sprintf("%s:%s", w.cfg.Server.Host, w.cfg.Server.Port),
		ReadTimeout:    time.Second * w.cfg.Server.ReadTimeout,
		WriteTimeout:   time.Second * w.cfg.Server.WriteTimeout,
		MaxHeaderBytes: maxHeaderBytes,
	}

	w.lock.Add(1)
	go func() {
		defer w.lock.Done()

		err := w.echo.StartServer(server)
		if err == nil || errors.Is(err, http.ErrServerClosed) {
			w.logger.Info("web server closed")
		} else {
			w.logger.Errorf("[ERROR] %+v", err)
		}
	}()

	return nil
}

func (w *Web) mountPoint() *echo.Group {
	return w.echo.Group("")
}

func (w *Web) Shutdown() error {
	c, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := w.echo.Shutdown(c); err != nil {
		return err
	}

	return nil
}
