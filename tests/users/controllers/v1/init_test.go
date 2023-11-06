package service

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"testing"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"github.com/redis/go-redis/v9"
)

const (
	dbDriver   = "postgres"
	dbSource   = "postgresql://test:test@localhost:5432/clean_db?sslmode=disable"
	migrations = "file://../../../../db/migrations"
	seed       = "file://../../../../db/seed"
)

// Declare a global variable to hold the Docker pool and resource.
var (
	network *dockertest.Network
)

func TestMain(m *testing.M) {
	// Initialize Docker pool
	// This command create a pool to interact with docker runtime
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not construct pool: %s", err)
	}

	// Ping the docker daemon
	// check if everything is good and
	// there is the connection with docker
	err = pool.Client.Ping()
	if err != nil {
		log.Fatalf("Could not connect to Docker: %s", err)
	}

	// Create a Docker network for the tests.
	network, err = pool.CreateNetwork("test-network")
	if err != nil {
		log.Fatalf("Could not create network: %v", err)
	}

	// Deploy the Postgres container.
	postgresResource, err := deployPostgres(pool)
	if err != nil {
		log.Fatalf("Could not start Postgres resource: %v", err)
	}

	// Deploy the Redis container.
	redisResource, err := deployRedis(pool)
	if err != nil {
		log.Fatalf("Could not start Redis resource: %v", err)
	}

	err = applyDatabaseMigrations()
	if err != nil {
		log.Fatalf("Could not apply postgres migration: %v", err)
	}

	err = applyDatabaseSeed()
	if err != nil {
		log.Fatalf("Could not apply postgres seed: %v", err)
	}

	// Deploy the API container.
	apiResource, err := deployAPIContainer(pool)
	if err != nil {
		log.Fatalf("Could not start clean service resource: %v", err)
	}

	resources := []*dockertest.Resource{
		postgresResource,
		redisResource,
		apiResource,
	}

	for _, res := range resources {
		// Kill container after 5 minite
		res.Expire(300)
	}

	// Run the tests.
	exitCode := m.Run()

	// Exit with the appropriate code.
	err = tearDown(pool, resources)
	if err != nil {
		log.Fatalf("Could not purge resource: %v", err)
	}

	os.Exit(exitCode)
}

func deployPostgres(pool *dockertest.Pool) (*dockertest.Resource, error) {
	resource, err := pool.RunWithOptions(&dockertest.RunOptions{
		Hostname:     "postgres",
		Name:         "clean-postgres",
		Repository:   "postgres",
		Tag:          "latest",
		ExposedPorts: []string{"5432"},
		PortBindings: map[docker.Port][]docker.PortBinding{
			"5432/tcp": {{HostIP: "", HostPort: "5432"}},
		},
		Env: []string{
			"POSTGRES_USER=test",
			"POSTGRES_PASSWORD=test",
			"POSTGRES_DB=clean_db",
			"listen_addresses = '*'",
		},
		Networks: []*dockertest.Network{
			network,
		},
	}, func(config *docker.HostConfig) {
		// set AutoRemove to true so that stopped container goes away by itself
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{
			Name: "no",
		}
	})
	if err != nil {
		return nil, fmt.Errorf("could not start resource: %v", err)
	}

	// Ensure the Postgres container is ready to accept connections.
	if err := pool.Retry(func() error {
		fmt.Println("Checking Postgres connection...")
		testDB, err := sqlx.Open(dbDriver, dbSource)
		if err != nil {
			return err
		}
		defer testDB.Close()

		return testDB.Ping()
	}); err != nil {
		return nil, fmt.Errorf("could not connect to docker: %v", err)
	}

	return resource, nil
}

// deployRedis builds and runs the Redis container.
func deployRedis(pool *dockertest.Pool) (*dockertest.Resource, error) {
	resource, err := pool.RunWithOptions(&dockertest.RunOptions{
		Hostname:     "redis",
		Name:         "clean-redis",
		Repository:   "redis",
		Tag:          "latest",
		ExposedPorts: []string{"6379"},
		PortBindings: map[docker.Port][]docker.PortBinding{
			"6379/tcp": {{HostIP: "", HostPort: "6379"}},
		},
		Networks: []*dockertest.Network{
			network,
		},
	}, func(config *docker.HostConfig) {
		// set AutoRemove to true so that stopped container goes away by itself
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{
			Name: "no",
		}
	})
	if err != nil {
		return nil, fmt.Errorf("could not start resource: %v", err)
	}

	// Ensure the Redis container is ready to accept connections.
	if err := pool.Retry(func() error {
		fmt.Println("Checking Redis connection...")
		db := redis.NewClient(&redis.Options{
			Addr:     "localhost:6379",
			Password: "",
			DB:       0,
		})

		_, err := db.Ping(context.Background()).Result()
		if err != nil {
			return err
		}

		defer db.Close()

		return nil
	}); err != nil {
		return nil, fmt.Errorf("could not connect to docker: %v", err)
	}

	return resource, nil
}

// deployAPIContainer builds and runs the API container.
func deployAPIContainer(pool *dockertest.Pool) (*dockertest.Resource, error) {
	// build and run the API container
	resource, err := pool.BuildAndRunWithBuildOptions(&dockertest.BuildOptions{
		ContextDir: "../../../..",
		Dockerfile: "deploy/Dockerfile.test",
	}, &dockertest.RunOptions{
		Hostname:     "clean",
		Name:         "api",
		ExposedPorts: []string{"8080"},
		PortBindings: map[docker.Port][]docker.PortBinding{
			"8080": {{HostIP: "0.0.0.0", HostPort: "8080"}},
		},
		Networks: []*dockertest.Network{
			network,
		},
	}, func(config *docker.HostConfig) {
		// set AutoRemove to true so that stopped container goes away by itself
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{
			Name: "no",
		}
	})
	if err != nil {
		return nil, fmt.Errorf("could not start resource: %v", err)
	}

	// check if the API container is ready to accept connections
	if err = pool.Retry(func() error {
		fmt.Println("Checking API connection...")
		_, err := http.Get("http://localhost:8080/healthz")
		if err != nil {
			return err
		}

		return nil
	}); err != nil {
		return nil, fmt.Errorf("could not start resource: %v", err)
	}

	return resource, nil
}

func applyDatabaseMigrations() error {
	fmt.Println("Apply Postgres migration...")
	testDB, err := sqlx.Open(dbDriver, dbSource)
	if err != nil {
		return err
	}
	defer testDB.Close()

	driver, err := postgres.WithInstance(testDB.DB, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("Could not create migration driver: %w", err)
	}

	migration, err := migrate.NewWithDatabaseInstance(
		migrations,
		"postgres",
		driver,
	)
	if err != nil {
		return fmt.Errorf("Failed to initialize migration instance: %w", err)
	}

	if err := migration.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("Failed to apply migrations: %w", err)
	}

	log.Println("Postgres migrations applied successfully")

	return nil
}

func applyDatabaseSeed() error {
	fmt.Println("Apply Postgres seed...")
	testDB, err := sqlx.Open(dbDriver, dbSource)
	if err != nil {
		return err
	}
	defer testDB.Close()

	driver, err := postgres.WithInstance(testDB.DB, &postgres.Config{
		MigrationsTable: "seed",
	})
	if err != nil {
		return fmt.Errorf("Could not create migration driver: %w", err)
	}

	migration, err := migrate.NewWithDatabaseInstance(
		seed,
		"postgres",
		driver,
	)
	if err != nil {
		return fmt.Errorf("Failed to initialize migration instance: %w", err)
	}

	if err := migration.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("Failed to apply postgres seed: %w", err)
	}

	log.Println("Database seed applied successfully")

	return nil
}

// tearDown purges the resources and removes the network.
func tearDown(pool *dockertest.Pool, resources []*dockertest.Resource) error {
	for _, resource := range resources {
		if err := pool.Purge(resource); err != nil {
			return fmt.Errorf("could not purge resource: %v", err)
		}
	}

	if err := pool.RemoveNetwork(network); err != nil {
		return fmt.Errorf("could not remove network: %v", err)
	}

	return nil
}
