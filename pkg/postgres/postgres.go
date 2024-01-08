package postgres

import (
	"context"
	"fmt"

	_ "github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"

	"github.com/KyKyPy3/clean/internal/infrastructure/config"
)

func New(ctx context.Context, cfg *config.PostgresConfig) (*sqlx.DB, error) {
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
