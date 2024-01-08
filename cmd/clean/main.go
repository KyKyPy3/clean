package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/KyKyPy3/clean/internal/infrastructure"
	"github.com/KyKyPy3/clean/internal/infrastructure/config"
	"github.com/KyKyPy3/clean/pkg/latch"
	"github.com/KyKyPy3/clean/pkg/logger"
	"github.com/KyKyPy3/clean/pkg/utils"
)

const (
	DefaultConfigFile = "config.yml"
	shutdownTimeout   = 30 * time.Second
)

func main() {
	lock := latch.NewCountDownLatch()

	ctx, cancel := context.WithCancel(context.Background())
	log.Println("Starting clean microservice")

	// Try to get config file name from environment
	configFile := utils.GetEnvVar("config", DefaultConfigFile)

	cfg, err := config.NewConfig(configFile)
	if err != nil {
		log.Fatalf("Loading config: %v", err)
	}

	// Init logger
	appLogger := logger.NewLogger(&cfg.Logger)
	appLogger.Init()
	appLogger.Infof("Version: %s, LogLevel: %s, SSL: %v", cfg.Server.Version, cfg.Logger.Level, cfg.Server.SSL)

	// Run our service
	srv := infrastructure.NewApp(ctx, cfg, appLogger, lock)
	if err = srv.Run(ctx); err != nil {
		log.Fatal(err)
	}

	// Listen for syscall signals for process to interrupt/quit
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-sig

	// Cancel context
	cancel()

	// Shutdown application
	srv.Shutdown()

	// Waiting for graceful shutdown
	lock.WaitWithTimeout(shutdownTimeout)
}
