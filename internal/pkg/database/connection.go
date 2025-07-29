package database

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/yantology/golang_template/internal/config"

	_ "github.com/lib/pq" // PostgreSQL driver
)

// Connect establishes a database connection using the provided configuration
func Connect(cfg config.DatabaseConfig) (*sql.DB, error) {
	// Open database connection
	db, err := sql.Open(cfg.GetDriverName(), cfg.GetDSN())
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	// Configure connection pool
	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	db.SetConnMaxLifetime(cfg.MaxLifetime)

	// Test the connection
	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return db, nil
}

// HealthCheck performs a health check on the database connection
func HealthCheck(db *sql.DB) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return fmt.Errorf("database health check failed: %w", err)
	}

	return nil
}

// GetDBStats returns database connection statistics
func GetDBStats(db *sql.DB) sql.DBStats {
	return db.Stats()
}