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

	// Initialize database
	db, err := database.NewDatabase(cfg)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to database")
	}

	// Initialize router with database
	r := router.SetupRouter(db)

	// Create Lambda adapter
	ginLambda = ginadapter.New(r)
}

func handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// Lambda will handle the request through Gin
	return ginLambda.Proxy(req)
}

func main() {
	lambda.Start(handler)
}
