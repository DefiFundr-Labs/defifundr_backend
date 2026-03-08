package authrouter

import (
	"time"

	"github.com/demola234/defifundr/infrastructure/common/logging"
	middleware "github.com/demola234/defifundr/infrastructure/middleware"
	authhandler "github.com/demola234/defifundr/internal/features/auth/handler"
	token "github.com/demola234/defifundr/pkg/token"
	"github.com/gin-gonic/gin"
)

// RegisterRoutes registers all authentication-related routes.
func RegisterRoutes(router *gin.Engine, handler *authhandler.Handler, tokenMaker token.Maker, logger logging.Logger) {
	// Apply global device tracking middleware
	router.Use(middleware.DeviceTrackingMiddleware())

	authRoutes := router.Group("/api/v1/auth")
	{
		authRoutes.Use(middleware.RateLimitMiddleware(5, time.Minute))

		authRoutes.POST("/web3auth/login", handler.Web3AuthLogin)
		authRoutes.POST("/register", handler.RegisterUser)
		authRoutes.POST("/login", handler.Login)
		authRoutes.POST("/refresh", handler.RefreshToken)
		authRoutes.POST("/forgot-password", handler.InitiatePasswordReset)
		authRoutes.POST("/verify-reset-otp", handler.VerifyResetOTP)
		authRoutes.POST("/reset-password", handler.ResetPassword)
	}

	authenticatedRoutes := router.Group("/api/v1/auth")
	authenticatedRoutes.Use(middleware.AuthMiddleware(tokenMaker, logger))
	{
		// Profile completion
		authenticatedRoutes.PUT("/profile/personal-details", handler.UpdatePersonalDetails)
		authenticatedRoutes.PUT("/profile/address", handler.UpdateAddressDetails)
		authenticatedRoutes.PUT("/profile/business", handler.UpdateBusinessDetails)
		authenticatedRoutes.GET("/profile/completion", handler.GetProfileCompletion)

		// Wallet management
		authenticatedRoutes.POST("/wallet/link", handler.LinkWallet)
		authenticatedRoutes.GET("/wallet", handler.GetWallets)

		// Device management
		authenticatedRoutes.GET("/security/devices", handler.GetUserDevices)
		authenticatedRoutes.POST("/security/devices/revoke", handler.RevokeDevice)

		// Security events
		authenticatedRoutes.GET("/security/events", handler.GetUserSecurityEvents)

		// MFA
		authenticatedRoutes.POST("/security/mfa/setup", handler.SetupMFA)
		authenticatedRoutes.POST("/security/mfa/verify", handler.VerifyMFA)

		// Session
		authenticatedRoutes.POST("/logout", handler.Logout)
	}
}
