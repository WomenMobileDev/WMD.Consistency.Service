package main

import (
	"context"

	"github.com/WomenMobileDev/WMD.Consistency.Service/internal/config"
	"github.com/WomenMobileDev/WMD.Consistency.Service/internal/database"
	"github.com/WomenMobileDev/WMD.Consistency.Service/internal/logger"
	"github.com/WomenMobileDev/WMD.Consistency.Service/internal/router"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	ginadapter "github.com/awslabs/aws-lambda-go-api-proxy/gin"
	"github.com/rs/zerolog/log"
)

var ginLambda *ginadapter.GinLambda

func init() {
	// Initialize logger
	logger.InitDefault()

	// Load configuration
	cfg := config.Load()
	logger.Init(cfg)

	log.Info().Msg("Lambda initialization starting...")
	log.Info().
		Str("db_host", cfg.Database.Host).
		Str("db_port", cfg.Database.Port).
		Str("db_name", cfg.Database.Name).
		Str("db_user", cfg.Database.User).
		Str("db_ssl_mode", cfg.Database.SSLMode).
		Msg("Database configuration loaded")

	// Initialize database with graceful handling
	log.Info().Msg("Attempting database connection...")
	db, err := database.NewDatabase(cfg)
	if err != nil {
		log.Warn().Err(err).Msg("Failed to connect to database, continuing without database connection")
		db = nil // Set to nil to indicate no database connection
	} else {
		log.Info().Msg("Database connected successfully!")
		// Run database migrations if connected
		if err := db.RunMigrations(); err != nil {
			log.Error().Err(err).Msg("Failed to run database migrations")
		} else {
			log.Info().Msg("Database migrations completed successfully")
		}
	}

	// Initialize router with database (can be nil)
	r := router.SetupRouter(db)

	// Create Lambda adapter
	ginLambda = ginadapter.New(r)

	log.Info().Msg("Lambda function initialized successfully")
}

func handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// Lambda will handle the request through Gin
	return ginLambda.Proxy(req)
}

func main() {
	lambda.Start(handler)
}
