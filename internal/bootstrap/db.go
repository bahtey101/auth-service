package bootstrap

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"auth-service/config"
)

func InitDB(ctx context.Context, cfg *config.Config) (*pgxpool.Pool, error) {
	pgxCfg, err := pgxpool.ParseConfig(cfg.PgDSN)
	if err != nil {
		return nil, err
	}

	pgxCfg.MaxConns = int32(cfg.PgMaxOpenConn)

	dbPool, err := pgxpool.NewWithConfig(ctx, pgxCfg)
	if err != nil {
		return nil, err
	}

	if err := dbPool.Ping(ctx); err != nil {
		return nil, err
	}

	return dbPool, nil
}
