package initialization

import (
	"fmt"
	"os"
	"strings"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"td27/rpc/basis/internal/config"
)

type writer struct {
	logger.Writer
}

// NewWriter writer constructor
func NewWriter(w logger.Writer) *writer {
	return &writer{Writer: w}
}

// Printf format and print SQL logs
func (w *writer) Printf(message string, data ...interface{}) {
	w.Writer.Printf(message, data...)
}

type gormLogWriter struct{}

func (gormLogWriter) Printf(format string, args ...interface{}) {
	fmt.Fprintf(os.Stdout, format, args...)
}

func gormConfig(logLevel string, logMode bool) *gorm.Config {
	var gormLogLevel logger.LogLevel
	switch strings.ToLower(logLevel) {
	case "silent":
		gormLogLevel = logger.Silent
	case "error":
		gormLogLevel = logger.Error
	case "warn":
		gormLogLevel = logger.Warn
	case "info":
		gormLogLevel = logger.Info
	default:
		gormLogLevel = logger.Warn
	}

	if !logMode {
		gormLogLevel = logger.Silent
	}

	newLogger := logger.New(
		NewWriter(gormLogWriter{}),
		logger.Config{
			SlowThreshold:             200 * time.Millisecond,
			LogLevel:                  gormLogLevel,
			IgnoreRecordNotFoundError: true,
			Colorful:                  false,
		},
	)

	return &gorm.Config{
		Logger:                                   newLogger,
		DisableForeignKeyConstraintWhenMigrating: true,
	}
}

// Gorm initializes PostgreSQL database connection
func Gorm(cfg config.Pgsql) (*gorm.DB, error) {
	if cfg.Dbname == "" {
		return nil, fmt.Errorf("database name is empty")
	}

	// PostgreSQL DSN format
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=Asia/Shanghai",
		cfg.Host, cfg.Username, cfg.Password, cfg.Dbname, cfg.Port)

	pgConfig := postgres.Config{
		DSN:                  dsn,
		PreferSimpleProtocol: false, // Enable prepared statement cache for better performance
	}

	db, err := gorm.Open(postgres.New(pgConfig), gormConfig(cfg.LogLevel, cfg.LogMode))
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(time.Hour)
	sqlDB.SetConnMaxIdleTime(30 * time.Minute)

	return db, nil
}

// IsAlreadyExistsError checks if error is PostgreSQL "already exists" condition
func IsAlreadyExistsError(err error) bool {
	if err == nil {
		return false
	}
	errStr := err.Error()
	return strings.Contains(errStr, "already exists") ||
		strings.Contains(errStr, "Duplicate key name") ||
		strings.Contains(errStr, "42P07") || // PostgreSQL: duplicate_table
		strings.Contains(errStr, "42710") // PostgreSQL: duplicate_object
}

// IsNotExistsError checks if error is PostgreSQL "does not exist" condition
func IsNotExistsError(err error) bool {
	if err == nil {
		return false
	}
	errStr := err.Error()
	return strings.Contains(errStr, "does not exist") ||
		strings.Contains(errStr, "42704") // PostgreSQL: undefined_object
}
