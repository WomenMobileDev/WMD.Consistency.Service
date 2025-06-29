package router

import (
	"net/http"

	"github.com/WomenMobileDev/WMD.Consistency.Service/internal/config"
	"github.com/WomenMobileDev/WMD.Consistency.Service/internal/database"
	"github.com/WomenMobileDev/WMD.Consistency.Service/internal/handlers"
	"github.com/WomenMobileDev/WMD.Consistency.Service/internal/middleware"
	"github.com/WomenMobileDev/WMD.Consistency.Service/internal/repository"
	"github.com/WomenMobileDev/WMD.Consistency.Service/internal/service"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// SetupRouter sets up the router with all routes and middleware
func SetupRouter(db *database.Database) *gin.Engine {
	r := gin.Default()

	// Add middlewares
	r.Use(cors.Default())
	r.Use(middleware.ErrorHandler())
	r.Use(middleware.ResponseFormatter())

	// Add NoRoute handler for proper 404 responses
	r.NoRoute(func(c *gin.Context) {
		middleware.RespondWithError(c, http.StatusNotFound, "NOT_FOUND", "The requested resource could not be found", gin.H{
			"documentation": "/swagger",
		})
	})

	// Load configuration
	cfg := config.Load()

	// Create repositories
	userRepo := repository.NewUserRepository(db.DB)
	habitRepo := repository.NewHabitRepository(db.DB)
	streakRepo := repository.NewStreakRepository(db.DB)
	checkInRepo := repository.NewCheckInRepository(db.DB)
	achievementRepo := repository.NewAchievementRepository(db.DB)

	// Create services
	authService := service.NewAuthService(userRepo, cfg)
	userService := service.NewUserService(userRepo, habitRepo, streakRepo, checkInRepo, achievementRepo)
	habitService := service.NewHabitService(habitRepo, streakRepo)
	streakService := service.NewStreakService(habitRepo, streakRepo, checkInRepo)
	checkInService := service.NewCheckInService(habitRepo, streakRepo, checkInRepo, achievementRepo)
	achievementService := service.NewAchievementService(achievementRepo, habitRepo)

	// Create handlers
	authHandler := handlers.NewAuthHandler(authService)
	userHandler := handlers.NewUserHandler(userService)
	habitHandler := handlers.NewHabitHandler(habitService)
	streakHandler := handlers.NewStreakHandler(streakService)
	checkInHandler := handlers.NewCheckInHandler(checkInService)
	achievementHandler := handlers.NewAchievementHandler(achievementService)

	// Root health check endpoints
	r.GET("/health", handlers.HealthCheck)
	r.GET("/db-health", handlers.DBHealthCheck(db))

	// Root welcome page
	r.GET("/", func(c *gin.Context) {
		middleware.RespondWithOK(c, gin.H{
			"name":          "Consistency API",
			"description":   "A RESTful API for tracking habits, streaks, and achievements",
			"version":       "1.0.0",
			"status":        "running",
			"documentation": "/swagger",
			"endpoints": map[string]string{
				"health":    "/health",
				"db_health": "/db-health",
				"api_v1":    "/api/v1",
			},
			"repository": "https://github.com/WomenMobileDev/WMD.Consistency.Service",
		})
	})

	// Setup Swagger documentation
	// Serve swagger.yaml file
	r.GET("/swagger.yaml", func(c *gin.Context) {
		c.File("./docs/swagger.yaml")
	})

	// Serve Swagger UI
	r.GET("/swagger", func(c *gin.Context) {
		swaggerHTML := `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>Consistency API - Swagger UI</title>
    <link rel="stylesheet" type="text/css" href="https://unpkg.com/swagger-ui-dist@4.5.0/swagger-ui.css" />
    <link rel="icon" type="image/png" href="https://unpkg.com/swagger-ui-dist@4.5.0/favicon-32x32.png" sizes="32x32" />
    <link rel="icon" type="image/png" href="https://unpkg.com/swagger-ui-dist@4.5.0/favicon-16x16.png" sizes="16x16" />
    <style>
        html { box-sizing: border-box; overflow: -moz-scrollbars-vertical; overflow-y: scroll; }
        *, *:before, *:after { box-sizing: inherit; }
        body { margin: 0; background: #fafafa; }
    </style>
</head>
<body>
    <div id="swagger-ui"></div>
    <script src="https://unpkg.com/swagger-ui-dist@4.5.0/swagger-ui-bundle.js" charset="UTF-8"></script>
    <script src="https://unpkg.com/swagger-ui-dist@4.5.0/swagger-ui-standalone-preset.js" charset="UTF-8"></script>
    <script>
        window.onload = function() {
            const ui = SwaggerUIBundle({
                url: "/swagger.yaml",
                dom_id: '#swagger-ui',
                deepLinking: true,
                presets: [
                    SwaggerUIBundle.presets.apis,
                    SwaggerUIStandalonePreset
                ],
                plugins: [
                    SwaggerUIBundle.plugins.DownloadUrl
                ],
                layout: "StandaloneLayout"
            });
            window.ui = ui;
        };
    </script>
</body>
</html>
`
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(swaggerHTML))
	})

	// API v1 routes
	v1 := r.Group("/api/v1")
	{
		// API v1 welcome page
		v1.GET("/", func(c *gin.Context) {
			middleware.RespondWithOK(c, gin.H{
				"name":        "Consistency API",
				"description": "A RESTful API for tracking habits, streaks, and achievements",
				"version":     "1.0.0",
				"status":      "running",
				"endpoints": map[string]string{
					"health":       "/api/v1/health",
					"db_health":    "/api/v1/db-health",
					"auth":         "/api/v1/auth",
					"profile":      "/api/v1/profile",
					"habits":       "/api/v1/habits",
					"achievements": "/api/v1/achievements",
				},
			})
		})
		// Health check endpoints
		v1.GET("/health", handlers.HealthCheck)
		v1.GET("/db-health", handlers.DBHealthCheck(db))

		// Auth routes
		auth := v1.Group("/auth")
		{
			auth.POST("/register", authHandler.Register())
			auth.POST("/login", authHandler.Login())
			auth.POST("/forgot-password", authHandler.ForgotPassword())
			auth.POST("/reset-password", authHandler.ResetPassword())
		}

		// Protected routes (require authentication)
		protected := v1.Group("")
		protected.Use(middleware.Auth())
		{
			// User routes
			protected.GET("/profile", userHandler.GetProfile())
			protected.PUT("/profile", userHandler.UpdateProfile())

			// Habit routes
			habits := protected.Group("/habits")
			{
				habits.POST("", habitHandler.CreateHabit())
				habits.GET("", habitHandler.ListHabits())
				habits.GET("/:id", habitHandler.GetHabit())
				habits.PUT("/:id", habitHandler.UpdateHabit())
				habits.DELETE("/:id", habitHandler.DeleteHabit())

				// Streak routes
				habits.POST("/:id/streaks", streakHandler.CreateStreak())
				habits.GET("/:id/streaks", streakHandler.ListStreaks())
				habits.GET("/:id/streaks/current", streakHandler.GetCurrentStreak())

				// Check-in routes
				habits.POST("/:id/check-ins", checkInHandler.CheckIn())
				habits.GET("/:id/check-ins", checkInHandler.ListCheckIns())

				// Achievement routes for a specific habit
				habits.GET("/:id/achievements", achievementHandler.ListHabitAchievements())
			}

			// Achievement routes
			achievements := protected.Group("/achievements")
			{
				achievements.GET("", achievementHandler.ListAchievements())
				achievements.GET("/:id", achievementHandler.GetAchievement())
			}
		}
	}

	return r
}
