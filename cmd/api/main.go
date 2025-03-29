package main

import (
	"context"
	"log"
	"time"

	"github.com/demola234/defifundr/cmd/api/docs"
	_ "github.com/demola234/defifundr/cmd/api/docs"
	"github.com/demola234/defifundr/config"
	db "github.com/demola234/defifundr/db/sqlc"
	"github.com/demola234/defifundr/infrastructure/common/logging"
	"github.com/demola234/defifundr/infrastructure/middleware"
	"github.com/demola234/defifundr/internal/adapters/handlers"
	"github.com/demola234/defifundr/internal/adapters/repositories"
	"github.com/demola234/defifundr/internal/adapters/routers"
	"github.com/demola234/defifundr/internal/core/services"
	tokenMaker "github.com/demola234/defifundr/pkg/token_maker"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title DefiFundr API
// @version 1.0
// @description Decentralized Payroll and Invoicing Platform for Remote Teams
// @termsOfService http://swagger.io/terms/
// @schemes http
// @contact.name DefiFundr Support
// @contact.url http://defifundr.com/support
// @contact.email hello@defifundr.com
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host localhost:8080
// @BasePath /api/v1
// @securityDefinitions.basic BasicAuth

func main() {
	// Load configuration
	configs, err := config.LoadConfig(".")
	if err != nil {
		log.Fatalf("cannot load config: %v", err)
	}
	// Initialize logger
	logger := logging.New(&configs)
	logger.Info("Starting application", map[string]interface{}{
		"environment": configs.Environment,
	})

	// Connect using pgx
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Connect to the database using the pgx driver with database/sql
	conn, err := pgxpool.New(ctx, "postgres://root:secret@localhost:5433/defi?sslmode=disable")
	if err != nil {
		logger.Fatal("Unable to connect to database", err, map[string]interface{}{
			"db_source": configs.DBSource,
		})
	}

	// Initialize repository
	dbQueries := db.New(conn)

	defer conn.Close()

	// Create repositories
	userRepo := repositories.NewUserRepository(*dbQueries)
	otpRepo := repositories.NewOtpRepository(*dbQueries)
	sessionRepo := repositories.NewSessionRepository(*dbQueries)
	waitlistRepo := repositories.NewWaitlistRepository(*dbQueries)
	
	tokenMaker, err := tokenMaker.NewTokenMaker(configs.TokenSymmetricKey)
	if err != nil {
		log.Fatalf("cannot create token maker: %v", err)
	}

	// Create email service
	emailService := services.NewEmailService(configs, logger)

	// Create services
	authService := services.NewAuthService(userRepo, otpRepo, sessionRepo, tokenMaker, configs)
	userService := services.NewUserService(userRepo)
	waitlistService := services.NewWaitlistService(waitlistRepo, emailService)

	// Create handlers
	authHandler := handlers.NewAuthHandler(authService, logger)
	userHandler := handlers.NewUserHandler(userService)
	waitlistHandler := handlers.NewWaitlistHandler(waitlistService, logger)

	// Initialize the router
	router := gin.New() // We use a custom logger middleware to log all requests

	// Apply our custom logging middleware
	router.Use(middleware.LoggingMiddleware(logger, &configs))
	router.Use(gin.Recovery()) // We still need recovery middleware

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Set up API routes
	setupRoutes(router, authHandler, userHandler, waitlistHandler, configs, logger)

	docs.SwaggerInfo.Title = "DefiFundr API"
	docs.SwaggerInfo.Description = "Decentralized Payroll and Invoicing Platform for Remote Teams"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = "localhost:8080"
	docs.SwaggerInfo.BasePath = "/api/v1"
	docs.SwaggerInfo.Schemes = []string{"http", "https"}

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Start the HTTP server
	logger.Info("HTTP server is running on", map[string]interface{}{
		"address": configs.HTTPServerAddress,
	})
	if err := router.Run(configs.HTTPServerAddress); err != nil {
		logger.Fatal("Failed to start HTTP server", err)
	}
}

// setupRoutes configures all the API routes
func setupRoutes(router *gin.Engine, authHandler *handlers.AuthHandler, userHandler *handlers.UserHandler, waitlistHandler *handlers.WaitlistHandler, configs config.Config, logger logging.Logger) {
	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "UP",
		})
	})

	v1 := router.Group("/api/v1")

	tokenMaker, err := tokenMaker.NewTokenMaker(configs.TokenSymmetricKey)
	if err != nil {
		logger.Panic("failed to create token maker", err)
	}

	// Middleware to check if the user is authenticated
	authMiddleware := middleware.AuthMiddleware(tokenMaker)

	// Register routes
	routers.RegisterAuthRoutes(v1, authHandler, authMiddleware)
	routers.RegisterUserRoutes(v1, userHandler, authMiddleware)
	routers.RegisterWaitlistRoutes(v1, waitlistHandler, authMiddleware)
}