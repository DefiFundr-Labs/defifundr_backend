package main

import (
	"context"
	"log"
	"time"

	"github.com/demola234/defifundr/cmd/api/docs"
	"github.com/demola234/defifundr/config"
	db "github.com/demola234/defifundr/db/sqlc"
	"github.com/demola234/defifundr/infrastructure/common/logging"
	"github.com/demola234/defifundr/infrastructure/mail"
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

// IMPORTANT: Do not include protocol in host (no http:// or https://)
// @title DefiFundr API
// @version 1.0
// @description Decentralized Payroll and Invoicing Platform for Remote Teams
// @termsOfService http://swagger.io/terms/
// @schemes http https
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
	conn, err := pgxpool.New(ctx, configs.DBSource)
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

	// Initialize Email System
	// Create AsyncQ email sender
	emailSender, err := mail.NewAsyncQEmailSender(configs, logger)
	if err != nil {
		logger.Fatal("Failed to create AsyncQ email sender", err, nil)
	}

	// Need to cast to access the non-interface methods
	asyncQSender, ok := emailSender.(*mail.AsyncQEmailSender)
	if !ok {
		logger.Fatal("Failed to cast email sender", nil, nil)
	}

	// Create the email worker with the async queue
	emailWorker, err := mail.NewEmailWorker(configs, logger, asyncQSender)
	if err != nil {
		logger.Fatal("Failed to create email worker", err, nil)
	}

	// Start the email worker
	emailWorker.Start()
	defer emailWorker.Stop()

	// Create email service using the email sender
	emailService := services.NewEmailService(configs, logger, emailSender)

	// Create services
	authService := services.NewAuthService(userRepo, otpRepo, sessionRepo, tokenMaker, configs)
	userService := services.NewUserService(userRepo)
	waitlistService := services.NewWaitlistService(waitlistRepo, emailService)

	// Create handlers
	authHandler := handlers.NewAuthHandler(authService, logger)
	userHandler := handlers.NewUserHandler(userService)
	waitlistHandler := handlers.NewWaitlistHandler(waitlistService, logger)

	// Initialize the router
	router := gin.New()

	// Apply our custom logging middleware
	router.Use(middleware.LoggingMiddleware(logger, &configs))
	router.Use(gin.Recovery())

	// Configure CORS to allow all origins
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

	// Explicitly set host based on environment without protocol
	var swaggerHost string
	if configs.Environment == "production" {
		swaggerHost = "defifundr.koyeb.app"
	} else {
		swaggerHost = "localhost:8080"
	}

	// Set Swagger info
	docs.SwaggerInfo.Title = "DefiFundr API"
	docs.SwaggerInfo.Description = "Decentralized Payroll and Invoicing Platform for Remote Teams"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = swaggerHost
	docs.SwaggerInfo.BasePath = "/api/v1"
	docs.SwaggerInfo.Schemes = []string{"https", "http"}

	// Setup Swagger endpoint
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
