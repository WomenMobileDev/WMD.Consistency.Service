package logger

import (
	"fmt"
	"os"
	"time"

	"github.com/WomenMobileDev/WMD.Consistency.Service/internal/config"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func Init(cfg *config.Config) {
	var output zerolog.ConsoleWriter
	if cfg.Logger.Pretty {
		output = zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: cfg.Logger.TimeFormat,
			NoColor:    false,
		}
	} else {
		output = zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: cfg.Logger.TimeFormat,
			NoColor:    true,
		}
	}

	log.Logger = zerolog.New(output).
		With().
		Timestamp().
		Caller().
		Logger()

	level, err := zerolog.ParseLevel(cfg.Logger.Level)
	if err != nil {
		fmt.Printf("Invalid log level '%s', defaulting to debug\n", cfg.Logger.Level)
		level = zerolog.DebugLevel
	}
	zerolog.SetGlobalLevel(level)

	log.Info().
		Str("level", level.String()).
		Bool("pretty", cfg.Logger.Pretty).
		Str("environment", cfg.Server.Env).
		Msg("Logger initialized")
}

func InitDefault() {
	output := zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: time.RFC3339,
		NoColor:    false,
	}

	log.Logger = zerolog.New(output).
		With().
		Timestamp().
		Caller().
		Logger()

	zerolog.SetGlobalLevel(zerolog.DebugLevel)

	log.Info().Msg("Default logger initialized")
}

func TestLogger() zerolog.Logger {
	output := zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: time.RFC3339,
		NoColor:    true,
	}

	return zerolog.New(output).
		With().
		Timestamp().
		Caller().
		Logger()
}
