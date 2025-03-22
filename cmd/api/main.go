package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/demola234/defifundr/internal/adapters/primary/api/middleware"
	db "github.com/demola234/defifundr/internal/adapters/secondary/db/postgres/sqlc"
	"github.com/demola234/defifundr/pkg/logging"
	"github.com/gin-gonic/gin"

	"github.com/demola234/defifundr/config"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib" // pgx driver for database/sql
)

func main() {
	// Load configuration
	configs, err := config.LoadConfig(".")
	if err != nil {
		log.Fatalf("cannot load config: %v", err)
	}

	// Initialize logger
	logger := logging.New(&configs)
	logger.Info("Starting application", map[string]interface{}{
		"environment": os.Getenv("ENVIRONMENT"),
		"version":     "1.0.0",
	})

	// Setup connection to database
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://postgres:postgres@localhost:5432/defifundr?sslmode=disable"
	}

	// Connect using pgx
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Connect to the database using the pgx driver with database/sql
	conn, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		logger.Fatal("Unable to connect to database", err, map[string]interface{}{
			"db_url": dbURL,
		})
	}
	defer conn.Close()

	logger.Info("Connected to database", map[string]interface{}{
		"db_url": dbURL,
	})

	// Initialize repository
	dbQueries := db.New(conn)

	// Set the gin mode based on environment
	if os.Getenv("ENVIRONMENT") != "development" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Initialize the router
	// We use a custom logger middleware to log all requests
	router := gin.New()

	// Apply our custom logging middleware
	router.Use(middleware.LoggingMiddleware(logger, &configs))
	router.Use(gin.Recovery()) // We still need recovery middleware

	// Set up API routes
	setupRoutes(router, dbQueries, logger)

	// Start the HTTP server
	logger.Info("HTTP server is running on", map[string]interface{}{
		"address": configs.HTTPServerAddress,
	})

	if err := router.Run(configs.HTTPServerAddress); err != nil {
		logger.Fatal("Failed to start HTTP server", err, nil)
	}
}

// setupRoutes configures all the API routes
func setupRoutes(router *gin.Engine, queries *db.Queries, logger logging.Logger) {
	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "UP",
		})
	})

	// API routes
	api := router.Group("/api/v1")
	{
		// User routes
		userRoutes := api.Group("/users")
		{
			userRoutes.POST("/", createUser(queries, logger))
			userRoutes.GET("/:id", getUser(queries, logger))
			// Add more user routes
		}

		// Wallet routes
		walletRoutes := api.Group("/wallets")
		{
			walletRoutes.POST("/", createWallet(queries, logger))
			walletRoutes.GET("/user/:userId", getUserWallets(queries, logger))
			// Add more wallet routes
		}

		// Organization routes
		orgRoutes := api.Group("/organizations")
		{
			orgRoutes.POST("/", createOrganization(queries, logger))
			orgRoutes.GET("/:id", getOrganization(queries, logger))
			// Add more organization routes
		}

		// Payroll routes
		payrollRoutes := api.Group("/payrolls")
		{
			payrollRoutes.POST("/", createPayroll(queries, logger))
			payrollRoutes.GET("/:id", getPayroll(queries, logger))
			// Add more payroll routes
		}

		// Invoice routes
		invoiceRoutes := api.Group("/invoices")
		{
			invoiceRoutes.POST("/", createInvoice(queries, logger))
			invoiceRoutes.GET("/:id", getInvoice(queries, logger))
			// Add more invoice routes
		}

		// Transaction routes
		txRoutes := api.Group("/transactions")
		{
			txRoutes.GET("/user/:userId", getUserTransactions(queries, logger))
			// Add more transaction routes
		}

		// Notification routes
		notificationRoutes := api.Group("/notifications")
		{
			notificationRoutes.GET("/user/:userId", getUserNotifications(queries, logger))
			// Add more notification routes
		}
	}
}

// Handler functions - these are placeholders that you'll need to implement
func createUser(queries *db.Queries, logger logging.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract request ID
		requestID, _ := c.Get("RequestID")
		reqLogger := logger.With("request_id", requestID)
		reqLogger.Debug("Processing create user request")

		// Implementation
		c.JSON(501, gin.H{"error": "Not implemented"})
	}
}

func getUser(queries *db.Queries, logger logging.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract request ID
		requestID, _ := c.Get("RequestID")
		reqLogger := logger.With("request_id", requestID)
		reqLogger.Debug("Processing get user request", map[string]interface{}{
			"user_id": c.Param("id"),
		})

		// Implementation
		c.JSON(501, gin.H{"error": "Not implemented"})
	}
}

func createWallet(queries *db.Queries, logger logging.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Implementation with logging
		requestID, _ := c.Get("RequestID")
		logger.With("request_id", requestID).Debug("Create wallet request")
		c.JSON(501, gin.H{"error": "Not implemented"})
	}
}

func getUserWallets(queries *db.Queries, logger logging.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Implementation with logging
		requestID, _ := c.Get("RequestID")
		logger.With("request_id", requestID).Debug("Get user wallets request")
		c.JSON(501, gin.H{"error": "Not implemented"})
	}
}

func createOrganization(queries *db.Queries, logger logging.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Implementation with logging
		requestID, _ := c.Get("RequestID")
		logger.With("request_id", requestID).Debug("Create organization request")
		c.JSON(501, gin.H{"error": "Not implemented"})
	}
}

func getOrganization(queries *db.Queries, logger logging.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Implementation with logging
		requestID, _ := c.Get("RequestID")
		logger.With("request_id", requestID).Debug("Get organization request")
		c.JSON(501, gin.H{"error": "Not implemented"})
	}
}

func createPayroll(queries *db.Queries, logger logging.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Implementation with logging
		requestID, _ := c.Get("RequestID")
		logger.With("request_id", requestID).Debug("Create payroll request")
		c.JSON(501, gin.H{"error": "Not implemented"})
	}
}

func getPayroll(queries *db.Queries, logger logging.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Implementation with logging
		requestID, _ := c.Get("RequestID")
		logger.With("request_id", requestID).Debug("Get payroll request")
		c.JSON(501, gin.H{"error": "Not implemented"})
	}
}

func createInvoice(queries *db.Queries, logger logging.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Implementation with logging
		requestID, _ := c.Get("RequestID")
		logger.With("request_id", requestID).Debug("Create invoice request")
		c.JSON(501, gin.H{"error": "Not implemented"})
	}
}

func getInvoice(queries *db.Queries, logger logging.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Implementation with logging
		requestID, _ := c.Get("RequestID")
		logger.With("request_id", requestID).Debug("Get invoice request")
		c.JSON(501, gin.H{"error": "Not implemented"})
	}
}

func getUserTransactions(queries *db.Queries, logger logging.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Implementation with logging
		requestID, _ := c.Get("RequestID")
		logger.With("request_id", requestID).Debug("Get user transactions request")
		c.JSON(501, gin.H{"error": "Not implemented"})
	}
}

func getUserNotifications(queries *db.Queries, logger logging.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Implementation with logging
		requestID, _ := c.Get("RequestID")
		logger.With("request_id", requestID).Debug("Get user notifications request")
		c.JSON(501, gin.H{"error": "Not implemented"})
	}
}
