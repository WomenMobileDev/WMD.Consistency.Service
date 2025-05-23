package config

import (
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
)

type Config struct {
	Server ServerConfig
	Logger LoggerConfig

	Database DatabaseConfig

	JWT JWTConfig
	Auth AuthConfig

	API APIConfig

	Monitoring MonitoringConfig
}

type ServerConfig struct {
	Port string
	Env  string
}
type LoggerConfig struct {
	Level      string
	Pretty     bool
	TimeFormat string
}

type DatabaseConfig struct {
	Host     string
	Port     string
	Name     string
	User     string
	Password string
	SSLMode  string
}

type JWTConfig struct {
	Secret string
	Expiry time.Duration
}

type AuthConfig struct {
	JWTSecret            string
	JWTExpiryHours      int
	PasswordResetExpiry time.Duration
	TokenIssuer         string
}
type APIConfig struct {
	Timeout   time.Duration
	RateLimit int
}
type MonitoringConfig struct {
	MetricsEnabled bool
	TracingEnabled bool
}

func Load() *Config {
	if err := godotenv.Load(); err != nil {
		log.Warn().Err(err).Msg("Error loading .env file, using environment variables")
	}

	return &Config{
		Server: ServerConfig{
			Port: getEnv("PORT", "8080"),
			Env:  getEnv("ENV", "development"),
		},
		Logger: LoggerConfig{
			Level:      getEnv("LOG_LEVEL", "debug"),
			Pretty:     getBoolEnv("LOG_PRETTY", true),
			TimeFormat: time.RFC3339,
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			Name:     getEnv("DB_NAME", "consistency_service"),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", "postgres"),
			SSLMode:  getEnv("DB_SSL_MODE", "disable"),
		},
		JWT: JWTConfig{
			Secret: getEnv("JWT_SECRET", "development_secret"),
			Expiry: getDurationEnv("JWT_EXPIRY", 24*time.Hour),
		},
		Auth: AuthConfig{
			JWTSecret:            getEnv("AUTH_JWT_SECRET", "habit_tracking_secret_key"),
			JWTExpiryHours:      getIntEnv("AUTH_JWT_EXPIRY_HOURS", 72),
			PasswordResetExpiry: getDurationEnv("AUTH_PASSWORD_RESET_EXPIRY", 24*time.Hour),
			TokenIssuer:         getEnv("AUTH_TOKEN_ISSUER", "habit-tracking-app"),
		},
		API: APIConfig{
			Timeout:   getDurationEnv("API_TIMEOUT", 30*time.Second),
			RateLimit: getIntEnv("API_RATE_LIMIT", 100),
		},
		Monitoring: MonitoringConfig{
			MetricsEnabled: getBoolEnv("METRICS_ENABLED", false),
			TracingEnabled: getBoolEnv("TRACING_ENABLED", false),
		},
	}
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return strings.TrimSpace(value)
}

func getBoolEnv(key string, defaultValue bool) bool {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}

	parsedValue, err := strconv.ParseBool(strings.TrimSpace(value))
	if err != nil {
		return defaultValue
	}

	return parsedValue
}

func getIntEnv(key string, defaultValue int) int {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}

	parsedValue, err := strconv.Atoi(strings.TrimSpace(value))
	if err != nil {
		return defaultValue
	}

	return parsedValue
}

func getDurationEnv(key string, defaultValue time.Duration) time.Duration {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}

	parsedValue, err := time.ParseDuration(strings.TrimSpace(value))
	if err != nil {
		return defaultValue
	}

	return parsedValue
}
