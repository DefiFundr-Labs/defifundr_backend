package main

import (
	"context"
	"log"
	"time"

	"github.com/demola234/defifundr/cmd/api/docs"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/demola234/defifundr/config"
	db "github.com/demola234/defifundr/db/sqlc"
	"github.com/demola234/defifundr/infrastructure/common/logging"
	"github.com/demola234/defifundr/infrastructure/mail"
	middlewareLocal "github.com/demola234/defifundr/infrastructure/middleware"
	authhandler "github.com/demola234/defifundr/internal/features/auth/handler"
	authrepo "github.com/demola234/defifundr/internal/features/auth/repository"
	authrouter "github.com/demola234/defifundr/internal/features/auth/router"
	authusecase "github.com/demola234/defifundr/internal/features/auth/usecase"
	userhandler "github.com/demola234/defifundr/internal/features/user/handler"
	userrepo "github.com/demola234/defifundr/internal/features/user/repository"
	userrouter "github.com/demola234/defifundr/internal/features/user/router"
	userusecase "github.com/demola234/defifundr/internal/features/user/usecase"
	waitlisthandler "github.com/demola234/defifundr/internal/features/waitlist/handler"
	waitlistrepo "github.com/demola234/defifundr/internal/features/waitlist/repository"
	waitlistrouter "github.com/demola234/defifundr/internal/features/waitlist/router"
	waitlistusecase "github.com/demola234/defifundr/internal/features/waitlist/usecase"
	adminhandler "github.com/demola234/defifundr/internal/features/admin/handler"
	adminrouter "github.com/demola234/defifundr/internal/features/admin/router"
	blockchainhandler "github.com/demola234/defifundr/internal/features/blockchain/handler"
	blockchainrouter "github.com/demola234/defifundr/internal/features/blockchain/router"
	companyhandler "github.com/demola234/defifundr/internal/features/company/handler"
	companyrouter "github.com/demola234/defifundr/internal/features/company/router"
	compliancehandler "github.com/demola234/defifundr/internal/features/compliance/handler"
	compliancerouter "github.com/demola234/defifundr/internal/features/compliance/router"
	contracthandler "github.com/demola234/defifundr/internal/features/contract/handler"
	contractrouter "github.com/demola234/defifundr/internal/features/contract/router"
	hrhandler "github.com/demola234/defifundr/internal/features/hr/handler"
	hrrouter "github.com/demola234/defifundr/internal/features/hr/router"
	invoicehandler "github.com/demola234/defifundr/internal/features/invoice/handler"
	invoicerouter "github.com/demola234/defifundr/internal/features/invoice/router"
	kychandler "github.com/demola234/defifundr/internal/features/kyc/handler"
	kycrouter "github.com/demola234/defifundr/internal/features/kyc/router"
	networkhandler "github.com/demola234/defifundr/internal/features/network/handler"
	networkrouter "github.com/demola234/defifundr/internal/features/network/router"
	notificationhandler "github.com/demola234/defifundr/internal/features/notification/handler"
	notificationrouter "github.com/demola234/defifundr/internal/features/notification/router"
	payrollhandler "github.com/demola234/defifundr/internal/features/payroll/handler"
	payrollrouter "github.com/demola234/defifundr/internal/features/payroll/router"
	taxhandler "github.com/demola234/defifundr/internal/features/tax/handler"
	taxrouter "github.com/demola234/defifundr/internal/features/tax/router"
	timesheethandler "github.com/demola234/defifundr/internal/features/timesheet/handler"
	timesheetrouter "github.com/demola234/defifundr/internal/features/timesheet/router"
	transactionhandler "github.com/demola234/defifundr/internal/features/transaction/handler"
	transactionrouter "github.com/demola234/defifundr/internal/features/transaction/router"
	wallethandler "github.com/demola234/defifundr/internal/features/wallet/handler"
	walletrouter "github.com/demola234/defifundr/internal/features/wallet/router"
	token "github.com/demola234/defifundr/pkg/token"
	"github.com/demola234/defifundr/pkg/tracing"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

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
// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and the JWT token.

var (
	// Build information - these would typically be set via ldflags during build
	version   = "1.0.0"
	commit    = "dev"
	buildTime = "unknown"
)

func main() {
	// Load configuration
	configs, err := config.LoadConfig(".")
	if err != nil {
		log.Fatalf("cannot load config: %v", err)
	}

	// Initialize logger
	logger := logging.New(&configs)
	logger.Info("Starting application", map[string]any{
		"environment": configs.Environment,
		"version":     version,
		"commit":      commit,
		"build_time":  buildTime,
	})

	// Connect using pgx
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	conn, err := pgxpool.New(ctx, configs.DBSource)
	if err != nil {
		logger.Fatal("Unable to connect to database", err, map[string]any{
			"db_source": configs.DBSource,
		})
	}
	defer conn.Close()

	dbQueries := db.New(conn)

	// ── Repositories ──────────────────────────────────────────────
	uRepo := userrepo.New(*dbQueries)
	oAuthRepo := authrepo.NewOAuthRepository(*dbQueries, logger)
	sessionRepo := authrepo.NewSessionRepository(*dbQueries, conn)
	walletRepo := authrepo.NewWalletRepository(*dbQueries, conn)
	securityRepo := authrepo.NewSecurityRepository(*dbQueries)
	otpRepo := authrepo.NewOTPRepository(*dbQueries)
	waitlistRepo := waitlistrepo.New(*dbQueries)

	maker, err := token.NewTokenMaker(configs.TokenSymmetricKey)
	if err != nil {
		log.Fatalf("cannot create token maker: %v", err)
	}

	// ── Email system ──────────────────────────────────────────────
	asyncQSender, err := mail.NewAsyncQEmailSender(configs, logger)
	if err != nil {
		logger.Fatal("Failed to create AsyncQ email sender", err, nil)
	}
	emailWorker, err := mail.NewEmailWorker(configs, logger, asyncQSender)
	if err != nil {
		logger.Fatal("Failed to create email worker", err, nil)
	}
	emailWorker.Start()
	defer emailWorker.Stop()

	emailService := mail.NewEmailService(configs, logger, asyncQSender)

	// ── Services ──────────────────────────────────────────────────
	userService := userusecase.New(uRepo)
	authService := authusecase.New(uRepo, sessionRepo, oAuthRepo, walletRepo, securityRepo, emailService, maker, configs, logger, otpRepo, userService)
	waitlistService := waitlistusecase.New(waitlistRepo, emailService)

	// ── Handlers ──────────────────────────────────────────────────
	aHandler := authhandler.New(authService, logger)
	uHandler := userhandler.New(userService)
	wHandler := waitlisthandler.New(waitlistService, logger)
	companyHandler := companyhandler.New()
	walletHandler := wallethandler.New()
	transactionHandler := transactionhandler.New()
	payrollHandler := payrollhandler.New()
	invoiceHandler := invoicehandler.New()
	contractHandler := contracthandler.New()
	kycHandler := kychandler.New()
	timesheetHandler := timesheethandler.New()
	hrHandler := hrhandler.New()
	notificationHandler := notificationhandler.New()
	complianceHandler := compliancehandler.New()
	taxHandler := taxhandler.New()
	networkHandler := networkhandler.New()
	blockchainHandler := blockchainhandler.New()
	adminHandler := adminhandler.New()

	// ── OpenTelemetry ─────────────────────────────────────────────
	tracingCfg := tracing.Config{
		ServiceName:       "defifundr-api",
		ServiceVersion:    version,
		Environment:       configs.Environment,
		UseStdoutExporter: configs.Environment != "production",
	}
	otelShutdown, err := tracing.SetupOTel(context.Background(), tracingCfg)
	if err != nil {
		logger.Fatal("Failed to set up OpenTelemetry", err, map[string]any{
			"service": tracingCfg.ServiceName,
		})
	}
	defer func() {
		if err := otelShutdown(context.Background()); err != nil {
			logger.Error("Failed to shutdown OpenTelemetry", err, nil)
		}
	}()

	// ── Router ────────────────────────────────────────────────────
	router := gin.New()
	router.Use(middlewareLocal.LoggingMiddleware(logger, &configs))
	router.Use(gin.Recovery())
	router.Use(otelgin.Middleware("defifundr-api"))
	router.Use(middlewareLocal.PrometheusMiddleware())
	router.Use(cors.New(cors.Config{
		AllowOrigins:     getAllowedOrigins(configs.Environment),
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	router.GET("/health", func(c *gin.Context) {
		delay := c.Query("delay")
		if delay != "" {
			if d, err := time.ParseDuration(delay + "ms"); err == nil {
				time.Sleep(d)
			}
		}
		c.JSON(200, gin.H{
			"status":    "healthy",
			"timestamp": time.Now(),
			"service":   "defifundr-api",
			"version":   version,
		})
	})
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))

	setupRoutes(router, aHandler, uHandler, wHandler,
		companyHandler, walletHandler, transactionHandler, payrollHandler,
		invoiceHandler, contractHandler, kycHandler, timesheetHandler,
		hrHandler, notificationHandler, complianceHandler, taxHandler,
		networkHandler, blockchainHandler, adminHandler,
		maker, configs, logger)

	// ── Swagger ───────────────────────────────────────────────────
	var swaggerHost string
	if configs.Environment == "production" {
		swaggerHost = "defifundr.koyeb.app"
	} else {
		swaggerHost = "localhost:8080"
	}
	docs.SwaggerInfo.Title = "DefiFundr API"
	docs.SwaggerInfo.Description = "Decentralized Payroll and Invoicing Platform for Remote Teams"
	docs.SwaggerInfo.Version = version
	docs.SwaggerInfo.Host = swaggerHost
	docs.SwaggerInfo.BasePath = "/api/v1"
	if configs.Environment == "production" {
		docs.SwaggerInfo.Schemes = []string{"https", "http"}
	} else {
		docs.SwaggerInfo.Schemes = []string{"http"}
	}
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	logger.Info("HTTP server is running on", map[string]any{
		"address": configs.HTTPServerAddress,
	})
	if err := router.Run(configs.HTTPServerAddress); err != nil {
		logger.Fatal("Failed to start HTTP server", err)
	}
}

func getAllowedOrigins(env string) []string {
	if env == "production" {
		return []string{"https://defifundr.com", "https://app.defifundr.com"}
	}
	return []string{"http://localhost:3000", "http://localhost:3001", "http://localhost:8080"}
}

// setupRoutes configures all the API routes.
func setupRoutes(
	router *gin.Engine,
	aHandler *authhandler.Handler,
	uHandler *userhandler.Handler,
	wHandler *waitlisthandler.Handler,
	companyHandler *companyhandler.Handler,
	walletHandler *wallethandler.Handler,
	transactionHandler *transactionhandler.Handler,
	payrollHandler *payrollhandler.Handler,
	invoiceHandler *invoicehandler.Handler,
	contractHandler *contracthandler.Handler,
	kycHandler *kychandler.Handler,
	timesheetHandler *timesheethandler.Handler,
	hrHandler *hrhandler.Handler,
	notificationHandler *notificationhandler.Handler,
	complianceHandler *compliancehandler.Handler,
	taxHandler *taxhandler.Handler,
	networkHandler *networkhandler.Handler,
	blockchainHandler *blockchainhandler.Handler,
	adminHandler *adminhandler.Handler,
	maker token.Maker,
	_ config.Config,
	logger logging.Logger,
) {
	
	v1 := router.Group("/api/v1")
	authMiddleware := middlewareLocal.AuthMiddleware(maker, logger)

	authrouter.RegisterRoutes(router, aHandler, maker, logger)
	userrouter.RegisterRoutes(v1, uHandler, authMiddleware)
	waitlistrouter.RegisterRoutes(v1, wHandler, authMiddleware)
	companyrouter.RegisterRoutes(v1, companyHandler, authMiddleware)
	walletrouter.RegisterRoutes(v1, walletHandler, authMiddleware)
	transactionrouter.RegisterRoutes(v1, transactionHandler, authMiddleware)
	payrollrouter.RegisterRoutes(v1, payrollHandler, authMiddleware)
	invoicerouter.RegisterRoutes(v1, invoiceHandler, authMiddleware)
	contractrouter.RegisterRoutes(v1, contractHandler, authMiddleware)
	kycrouter.RegisterRoutes(v1, kycHandler, authMiddleware)
	timesheetrouter.RegisterRoutes(v1, timesheetHandler, authMiddleware)
	hrrouter.RegisterRoutes(v1, hrHandler, authMiddleware)
	notificationrouter.RegisterRoutes(v1, notificationHandler, authMiddleware)
	compliancerouter.RegisterRoutes(v1, complianceHandler, authMiddleware)
	taxrouter.RegisterRoutes(v1, taxHandler, authMiddleware)
	networkrouter.RegisterRoutes(v1, networkHandler, authMiddleware)
	blockchainrouter.RegisterRoutes(v1, blockchainHandler, authMiddleware)
	adminrouter.RegisterRoutes(v1, adminHandler, authMiddleware)
}
