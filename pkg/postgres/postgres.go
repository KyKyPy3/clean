package postgres

import (
	"context"
	"fmt"
	"time"

	_ "github.com/jackc/pgx/stdlib" // use in lib
	"github.com/jmoiron/sqlx"
)

const maxIdleConnections = 30

type Config struct {
	Host         string
	Port         string
	User         string
	Password     string
	DBName       string
	SSLMode      bool
	MaxOpenConn  int
	ConnLifetime time.Duration
	MaxIdleTime  time.Duration
}

func New(ctx context.Context, cfg Config) (*sqlx.DB, error) {
	dataSourceName := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s",
		cfg.Host,
		cfg.Port,
		cfg.User,
		cfg.DBName,
		cfg.Password,
	)

	db, err := sqlx.ConnectContext(ctx, "pgx", dataSourceName)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(cfg.MaxOpenConn)
	db.SetConnMaxLifetime(cfg.ConnLifetime)
	db.SetMaxIdleConns(maxIdleConnections)
	db.SetConnMaxIdleTime(cfg.MaxIdleTime)
	if err = db.PingContext(ctx); err != nil {
		return nil, err
	}

	return db, nil
}
