package postgres

import (
	"context"
	"fmt"
	"time"

	_ "github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"
)

type Config struct {
	Host         string
	Port         string
	User         string
	Password     string
	DbName       string
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
		cfg.DbName,
		cfg.Password,
	)

	db, err := sqlx.ConnectContext(ctx, "pgx", dataSourceName)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(cfg.MaxOpenConn)
	db.SetConnMaxLifetime(cfg.ConnLifetime)
	db.SetMaxIdleConns(30)
	db.SetConnMaxIdleTime(cfg.MaxIdleTime)
	if err = db.PingContext(ctx); err != nil {
		return nil, err
	}

	return db, nil
}
