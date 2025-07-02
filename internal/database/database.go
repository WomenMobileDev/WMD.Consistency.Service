package database

import (
	"context"
	"fmt"
	"time"

	"github.com/WomenMobileDev/WMD.Consistency.Service/internal/config"
	"github.com/rs/zerolog/log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Database struct {
	DB *gorm.DB
}

// NewDatabase creates a new GORM database connection
func NewDatabase(cfg *config.Config) (*Database, error) {
	// Build DSN (Data Source Name)
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s connect_timeout=%d",
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Name,
		cfg.Database.SSLMode,
		5, // 5 second connection timeout
	)

	// Configure GORM logger
	gormLogger := logger.New(
		&GormLogWriter{}, // Custom log writer that uses zerolog
		logger.Config{
			SlowThreshold:             200 * time.Millisecond,
			LogLevel:                  logger.Info,
			IgnoreRecordNotFoundError: true,
			Colorful:                  false,
		},
	)

	// Configure GORM with retry logic - reduced for Lambda
	var db *gorm.DB
	var err error
	var retryCount int
	maxRetries := 2               // Reduced from 5 to 2 for Lambda
	retryDelay := 1 * time.Second // Reduced from 2 seconds

	for retryCount < maxRetries {
		// Create a context with timeout for each connection attempt
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

		// Open GORM connection
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
			Logger: gormLogger,
			NowFunc: func() time.Time {
				return time.Now().UTC()
			},
		})

		if err == nil {
			// Test connection with timeout
			sqlDB, dbErr := db.DB()
			if dbErr == nil {
				err = sqlDB.PingContext(ctx)
				if err == nil {
					// Configure connection pool settings
					sqlDB.SetMaxIdleConns(5)                   // Reduced for Lambda
					sqlDB.SetMaxOpenConns(10)                  // Reduced for Lambda
					sqlDB.SetConnMaxLifetime(30 * time.Minute) // Shorter for Lambda
					cancel()
					break
				}
			} else {
				err = dbErr
			}
		}

		cancel()
		retryCount++
		log.Warn().
			Err(err).
			Int("retry", retryCount).
			Int("maxRetries", maxRetries).
			Msg("Failed to connect to database, retrying...")

		if retryCount < maxRetries {
			time.Sleep(retryDelay)
		}
	}

	if err != nil {
		return nil, fmt.Errorf("failed to connect to database after %d attempts: %w", maxRetries, err)
	}

	log.Info().
		Str("host", cfg.Database.Host).
		Str("port", cfg.Database.Port).
		Str("database", cfg.Database.Name).
		Str("user", cfg.Database.User).
		Msg("Successfully connected to database")

	return &Database{DB: db}, nil
}

// Close closes the database connection
func (db *Database) Close() {
	if db.DB != nil {
		sqlDB, err := db.DB.DB()
		if err == nil {
			err = sqlDB.Close()
			if err != nil {
				log.Error().Err(err).Msg("Error closing database connection")
			} else {
				log.Info().Msg("Database connection closed")
			}
		}
	}
}

// Ping checks if the database connection is alive
func (db *Database) Ping(ctx context.Context) error {
	sqlDB, err := db.DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.PingContext(ctx)
}

// GormLogWriter is a custom log writer for GORM that uses zerolog
type GormLogWriter struct{}

// Printf implements the logger.Writer interface for GORM
func (w *GormLogWriter) Printf(format string, args ...interface{}) {
	log.Debug().Msgf(format, args...)
}

// AutoMigrate runs migrations for the provided models
func (db *Database) AutoMigrate(models ...interface{}) error {
	return db.DB.AutoMigrate(models...)
}

// CreateRecord creates a new record in the database
func (db *Database) CreateRecord(value interface{}) error {
	return db.DB.Create(value).Error
}

// FindRecord finds a record by its primary key
func (db *Database) FindRecord(dest interface{}, primaryKey interface{}) error {
	return db.DB.First(dest, primaryKey).Error
}

// UpdateRecord updates a record in the database
func (db *Database) UpdateRecord(value interface{}) error {
	return db.DB.Save(value).Error
}

// DeleteRecord deletes a record from the database
func (db *Database) DeleteRecord(value interface{}) error {
	return db.DB.Delete(value).Error
}
