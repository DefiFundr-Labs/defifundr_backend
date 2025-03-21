package main

import (
	"context"
	"log"
	"os"
	"time"

	db "github.com/defifundr/defifundr-backend/internal/adapters/secondary/db/postgres/sqlc"
	"github.com/gin-gonic/gin"

	"github.com/defifundr/defifundr-backend/config"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib" // pgx driver for database/sql
)

func main() {
	// Load configuration
	configs, err := config.LoadConfig(".")
	if err != nil {
		log.Fatalf("cannot load config: %v", err)
	}

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
		log.Fatalf("Unable to connect to database: %v", err)
	}
	defer conn.Close()

	// Initialize repository
	dbQueries := db.New(conn)

	// Initialize the router
	router := gin.Default()

	// Set up API routes
	setupRoutes(router, dbQueries)

	// Start the HTTP server
	log.Printf("HTTP server is running on %s", configs.HTTPServerAddress)
	if err := router.Run(configs.HTTPServerAddress); err != nil {
		log.Fatalf("Failed to start HTTP server: %v", err)
	}
}

// setupRoutes configures all the API routes
func setupRoutes(router *gin.Engine, queries *db.Queries) {
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
			userRoutes.POST("/", createUser(queries))
			userRoutes.GET("/:id", getUser(queries))
			// Add more user routes
		}

		// Wallet routes
		walletRoutes := api.Group("/wallets")
		{
			walletRoutes.POST("/", createWallet(queries))
			walletRoutes.GET("/user/:userId", getUserWallets(queries))
			// Add more wallet routes
		}

		// Organization routes
		orgRoutes := api.Group("/organizations")
		{
			orgRoutes.POST("/", createOrganization(queries))
			orgRoutes.GET("/:id", getOrganization(queries))
			// Add more organization routes
		}

		// Payroll routes
		payrollRoutes := api.Group("/payrolls")
		{
			payrollRoutes.POST("/", createPayroll(queries))
			payrollRoutes.GET("/:id", getPayroll(queries))
			// Add more payroll routes
		}

		// Invoice routes
		invoiceRoutes := api.Group("/invoices")
		{
			invoiceRoutes.POST("/", createInvoice(queries))
			invoiceRoutes.GET("/:id", getInvoice(queries))
			// Add more invoice routes
		}

		// Transaction routes
		txRoutes := api.Group("/transactions")
		{
			txRoutes.GET("/user/:userId", getUserTransactions(queries))
			// Add more transaction routes
		}

		// Notification routes
		notificationRoutes := api.Group("/notifications")
		{
			notificationRoutes.GET("/user/:userId", getUserNotifications(queries))
			// Add more notification routes
		}
	}
}

// Handler functions - these are placeholders that you'll need to implement
func createUser(queries *db.Queries) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Implementation
		c.JSON(501, gin.H{"error": "Not implemented"})
	}
}

func getUser(queries *db.Queries) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Implementation
		c.JSON(501, gin.H{"error": "Not implemented"})
	}
}

func createWallet(queries *db.Queries) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Implementation
		c.JSON(501, gin.H{"error": "Not implemented"})
	}
}

func getUserWallets(queries *db.Queries) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Implementation
		c.JSON(501, gin.H{"error": "Not implemented"})
	}
}

func createOrganization(queries *db.Queries) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Implementation
		c.JSON(501, gin.H{"error": "Not implemented"})
	}
}

func getOrganization(queries *db.Queries) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Implementation
		c.JSON(501, gin.H{"error": "Not implemented"})
	}
}

func createPayroll(queries *db.Queries) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Implementation
		c.JSON(501, gin.H{"error": "Not implemented"})
	}
}

func getPayroll(queries *db.Queries) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Implementation
		c.JSON(501, gin.H{"error": "Not implemented"})
	}
}

func createInvoice(queries *db.Queries) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Implementation
		c.JSON(501, gin.H{"error": "Not implemented"})
	}
}

func getInvoice(queries *db.Queries) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Implementation
		c.JSON(501, gin.H{"error": "Not implemented"})
	}
}

func getUserTransactions(queries *db.Queries) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Implementation
		c.JSON(501, gin.H{"error": "Not implemented"})
	}
}

func getUserNotifications(queries *db.Queries) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Implementation
		c.JSON(501, gin.H{"error": "Not implemented"})
	}
}
