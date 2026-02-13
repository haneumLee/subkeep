package database

import (
	"database/sql"
	"log/slog"

	"github.com/subkeep/backend/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DB holds the global database connection instance.
var DB *gorm.DB

// Connect initializes the PostgreSQL connection using GORM.
func Connect(cfg *config.DatabaseConfig) (*gorm.DB, error) {
	gormLogLevel := logger.Warn
	if cfg.SSLMode == "disable" {
		// Likely development; use more verbose logging.
		gormLogLevel = logger.Info
	}

	db, err := gorm.Open(postgres.Open(cfg.DSN()), &gorm.Config{
		Logger:                 logger.Default.LogMode(gormLogLevel),
		SkipDefaultTransaction: true,
		PrepareStmt:            true,
	})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	// Configure connection pool.
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(cfg.ConnMaxLifetime)

	if err := sqlDB.Ping(); err != nil {
		return nil, err
	}

	DB = db
	slog.Info("database connected",
		"host", cfg.Host,
		"port", cfg.Port,
		"database", cfg.Name,
		"max_open_conns", cfg.MaxOpenConns,
		"max_idle_conns", cfg.MaxIdleConns,
		"conn_max_lifetime", cfg.ConnMaxLifetime,
	)

	return db, nil
}

// Close gracefully closes the database connection.
func Close() error {
	if DB == nil {
		return nil
	}

	sqlDB, err := DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

// HealthCheck verifies the database connection is alive.
func HealthCheck() error {
	if DB == nil {
		return sql.ErrConnDone
	}

	sqlDB, err := DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Ping()
}
