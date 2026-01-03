package pg

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/augustdev/autoclip/internal/storage/pg/generated/users"
	pgxdecimal "github.com/jackc/pgx-shopspring-decimal"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DbConfig struct {
	Url             string        `mapstructure:"url"`
	MaxConns        int32         `mapstructure:"maxConns"`
	MinConns        int32         `mapstructure:"minConns"`
	MaxConnLifetime time.Duration `mapstructure:"maxConnLifetime"`
	MaxConnIdleTime time.Duration `mapstructure:"maxConnIdleTime"`
}

type DB struct {
	*pgxpool.Pool
	logger      *slog.Logger
	userQueries users.Querier
}

func NewDatabase(config DbConfig, logger *slog.Logger) (*DB, error) {
	ctx := context.Background()

	poolConfig, err := pgxpool.ParseConfig(config.Url)
	if err != nil {
		return nil, fmt.Errorf("failed to parse database config: %w", err)
	}

	if config.MaxConns > 0 {
		poolConfig.MaxConns = config.MaxConns
	} else {
		poolConfig.MaxConns = 10
	}

	if config.MinConns > 0 {
		poolConfig.MinConns = config.MinConns
	} else {
		poolConfig.MinConns = 2
	}

	if config.MaxConnLifetime > 0 {
		poolConfig.MaxConnLifetime = config.MaxConnLifetime
	} else {
		poolConfig.MaxConnLifetime = 30 * time.Minute
	}

	if config.MaxConnIdleTime > 0 {
		poolConfig.MaxConnIdleTime = config.MaxConnIdleTime
	} else {
		poolConfig.MaxConnIdleTime = 5 * time.Minute
	}

	logger.Info("Database pool configuration",
		"maxConns", poolConfig.MaxConns,
		"minConns", poolConfig.MinConns,
		"maxConnLifetime", poolConfig.MaxConnLifetime,
		"maxConnIdleTime", poolConfig.MaxConnIdleTime,
	)

	poolConfig.AfterConnect = func(ctx context.Context, conn *pgx.Conn) error {
		pgxdecimal.Register(conn.TypeMap())
		return nil
	}

	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection pool: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	logger.Info("Successfully connected to database")

	logger.Info("Running database migrations...")
	if err := RunMigrations(pool); err != nil {
		logger.Error("Failed to run database migrations", "error", err)
		return nil, fmt.Errorf("failed to run database migrations: %w", err)
	}
	logger.Info("Database migrations completed successfully")

	return &DB{
		Pool:        pool,
		logger:      logger,
		userQueries: users.New(pool),
	}, nil
}

func (db *DB) Health() error {
	return db.Ping(context.Background())
}

func NewUserQueries(database *DB) users.Querier {
	return database.userQueries
}
