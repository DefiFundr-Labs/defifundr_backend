package routers

import (
    "time"

    "github.com/demola234/defifundr/infrastructure/common/logging"
    middleware "github.com/demola234/defifundr/infrastructure/middleware"
    "github.com/demola234/defifundr/internal/adapters/handlers"
    token_maker "github.com/demola234/defifundr/pkg/token_maker"
    "github.com/gin-gonic/gin"
)

// RegisterAuthRoutes registers all authentication related routes
func RegisterAuthRoutes(router *gin.Engine, authHandler *handlers.AuthHandler, tokenMaker token_maker.Maker, logger logging.Logger) {
    // Apply global device tracking middleware
    router.Use(middleware.DeviceTrackingMiddleware())

    // Auth routes group
    authRoutes := router.Group("/api/v1/auth")
    {
        // Rate limit login and registration routes
        authRoutes.Use(middleware.RateLimitMiddleware(5, time.Minute))

        // Web3Auth login/registration
        authRoutes.POST("/web3auth/login", authHandler.Web3AuthLogin)

        // Email authentication
        authRoutes.POST("/register", authHandler.RegisterUser)
        authRoutes.POST("/login", authHandler.Login)
        authRoutes.POST("/refresh", authHandler.RefreshToken)

        // Forgot password
        authRoutes.POST("/forgot-password", authHandler.InitiatePasswordReset)
        authRoutes.POST("/verify-reset-otp", authHandler.VerifyResetOTP)
        authRoutes.POST("/reset-password", authHandler.ResetPassword)
    }

    // Authenticated routes
    authenticatedRoutes := router.Group("/api/v1/auth")
    authenticatedRoutes.Use(middleware.AuthMiddleware(tokenMaker, logger))
    {
        // User profile completion
        authenticatedRoutes.PUT("/profile/personal-details", authHandler.UpdatePersonalDetails)
        authenticatedRoutes.PUT("/profile/address", authHandler.UpdateAddressDetails)
        authenticatedRoutes.PUT("/profile/business", authHandler.UpdateBusinessDetails)
        authenticatedRoutes.GET("/profile/completion", authHandler.GetProfileCompletion)

        // Wallet management
        authenticatedRoutes.POST("/wallet/link", authHandler.LinkWallet)
        authenticatedRoutes.GET("/wallet", authHandler.GetWallets)

        // Device management
        authenticatedRoutes.GET("/security/devices", authHandler.GetUserDevices)
        authenticatedRoutes.POST("/security/devices/revoke", authHandler.RevokeDevice)

        // Security events
        authenticatedRoutes.GET("/security/events", authHandler.GetUserSecurityEvents)

        // MFA
        authenticatedRoutes.POST("/security/mfa/setup", authHandler.SetupMFA)
        authenticatedRoutes.POST("/security/mfa/verify", authHandler.VerifyMFA)

        // Session management
        authenticatedRoutes.POST("/logout", authHandler.Logout)

    }

    // High-security routes requiring MFA
    secureRoutes := router.Group("/api/v1/secure")
    secureRoutes.Use(middleware.AuthMiddleware(tokenMaker, logger))
    secureRoutes.Use(middleware.MFARequiredMiddleware(authHandler.GetUserRepository()))
    {
        // Add high-security routes here, like wallet withdrawals or changes to account settings
    }
}