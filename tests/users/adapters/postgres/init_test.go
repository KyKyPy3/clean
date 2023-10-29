package postgres_test

import (
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/KyKyPy3/clean/config"
	psql "github.com/KyKyPy3/clean/internal/user/adapters/postgres"
	"github.com/KyKyPy3/clean/internal/user/domain/service"
	"github.com/KyKyPy3/clean/pkg/logger"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
)

const (
	dbDriver   = "postgres"
	dbSource   = "postgresql://test:test@localhost:%s/clean_db?sslmode=disable"
	migrations = "file://../../../../db/migrations"
	seed       = "file://../../../../db/seed"
)

var (
	testDB   *sqlx.DB
	repo     service.UserPgStorage
	pool     *dockertest.Pool
	resource *dockertest.Resource
)

func TestMain(m *testing.M) {
	// Set up the Docker test environment
	setupDockerTestEnvironment()

	// Apply database migrations
	applyDatabaseMigrations()

	// Apply database seed
	applyDatabaseSeed()

	// Create logger
	// TODO: add discard logger here
	loggerCfg := &config.LoggerConfig{Mode: "test"}
	logger := logger.NewLogger(loggerCfg)
	logger.Init()

	repo = psql.NewUserPgStorage(testDB, logger)

	// Run the tests
	code := m.Run()

	// Tear down the Docker test environment
	tearDownDockerTestEnvironment()

	// Exit with the appropriate exit code
	os.Exit(code)
}

func setupDockerTestEnvironment() {
	setupTimeoutDuration := 5 * time.Minute
	setupDone := make(chan bool)

	go func() {
		var err error
		pool, err = dockertest.NewPool("")
		if err != nil {
			log.Fatalf("Could not construct pool: %s", err)
		}

		err = pool.Client.Ping()
		if err != nil {
			log.Fatalf("Could not connect to Docker: %s", err)
		}

		resource, err = pool.RunWithOptions(&dockertest.RunOptions{
			Repository: "postgres",
			Tag:        "latest",
			Env: []string{
				"POSTGRES_USER=test",
				"POSTGRES_PASSWORD=test",
				"POSTGRES_DB=clean_db",
				"listen_addresses = '*'",
			},
		}, func(config *docker.HostConfig) {
			config.AutoRemove = true
			config.RestartPolicy = docker.RestartPolicy{
				Name: "no",
			}
		})

		if err != nil {
			log.Fatalf("Failed to create resource: %s", err)
		}

		if err := pool.Retry(func() error {
			var err error
			testDB, err = sqlx.Open(dbDriver, fmt.Sprintf(dbSource, resource.GetPort("5432/tcp")))
			if err != nil {
				return err
			}

			return testDB.Ping()
		}); err != nil {
			log.Fatalf("Could not connect to database: %s", err)
		}

		setupDone <- true
	}()

	select {
	case <-setupDone:
		log.Println("Docker test environment setup completed")
	case <-time.After(setupTimeoutDuration):
		log.Println("Docker test environment setup timed out")
	}
}

func applyDatabaseMigrations() {
	driver, err := postgres.WithInstance(testDB.DB, &postgres.Config{})
	if err != nil {
		log.Fatalf("Could not create migration driver: %s", err)
	}

	migration, err := migrate.NewWithDatabaseInstance(
		migrations,
		"postgres",
		driver,
	)
	if err != nil {
		log.Fatalf("Failed to initialize migration instance: %s", err)
	}

	if err := migration.Up(); err != nil && err != migrate.ErrNoChange {
		closeMigrate(migration)

		log.Fatalf("Failed to apply migrations: %s", err)
	}

	log.Println("Database migrations applied successfully")
}

func applyDatabaseSeed() {
	driver, err := postgres.WithInstance(testDB.DB, &postgres.Config{
		MigrationsTable: "seed",
	})
	if err != nil {
		log.Fatalf("Could not create migration driver: %s", err)
	}

	migration, err := migrate.NewWithDatabaseInstance(
		seed,
		"postgres",
		driver,
	)
	if err != nil {
		log.Fatalf("Failed to initialize migration instance: %s", err)
	}

	if err := migration.Up(); err != nil && err != migrate.ErrNoChange {
		closeMigrate(migration)
		log.Fatalf("Failed to apply seed: %s", err)
	}

	log.Println("Database seed applied successfully")
}

func tearDownDockerTestEnvironment() {
	if err := pool.Purge(resource); err != nil {
		log.Fatalf("Could not purge Docker resource: %s", err)
	}
}

func closeMigrate(migrate *migrate.Migrate) {
	sourceErr, databaseErr := migrate.Close()
	if sourceErr != nil {
		log.Fatal("error closing migration source", sourceErr)
	}
	if databaseErr != nil {
		log.Fatal("error closing database source", databaseErr)
	}
}
