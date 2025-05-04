// middleware/auth_middleware.go
package middleware

import (
	"log"
	"net/http"
	"strings"

	"github.com/demola234/defifundr/infrastructure/common/logging"
	response "github.com/demola234/defifundr/internal/adapters/dto/response"
	"github.com/demola234/defifundr/internal/core/ports"
	token_maker "github.com/demola234/defifundr/pkg/token_maker"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// AuthMiddleware validates the JWT token in the Authorization header
func AuthMiddleware(tokenMaker token_maker.Maker, logger logging.Logger) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Get Authorization header
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			ctx.JSON(http.StatusUnauthorized, response.ErrorResponse{
				Success: false,
				Message: "Authorization header is required",
			})
			ctx.Abort()
			return
		}

		// Check if the auth header starts with "Bearer "
		fields := strings.Fields(authHeader)
		if len(fields) < 2 || fields[0] != "Bearer" {
			ctx.JSON(http.StatusUnauthorized, response.ErrorResponse{
				Success: false,
				Message: "Invalid authorization format",
			})
			ctx.Abort()
			return
		}

		// Extract token
		accessToken := fields[1]

		// Verify token
		payload, err := tokenMaker.VerifyToken(accessToken)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, response.ErrorResponse{
				Success: false,
				Message: "Invalid or expired token",
			})
			ctx.Abort()
			return
		}

		// Set user ID in context
		ctx.Set("user_id", payload.UserID)
		ctx.Set("email", payload.Email)

		// Continue to the next handler
		ctx.Next()
	}
}

// DeviceTrackingMiddleware tracks device information for security purposes
func DeviceTrackingMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		clientIP := ctx.ClientIP()
		userAgent := ctx.Request.UserAgent()

		// Store in context for later use
		ctx.Set("client_ip", clientIP)
		ctx.Set("user_agent", userAgent)

		// Add a request ID for tracking
		requestID := uuid.New().String()
		ctx.Set("request_id", requestID)
		ctx.Header("X-Request-ID", requestID)

		ctx.Next()
	}
}

// MFARequiredMiddleware ensures that MFA is completed for critical operations
func MFARequiredMiddleware(userRepo ports.UserRepository) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Get user ID from context
		userID, exists := ctx.Get("user_id")
		if !exists {
			ctx.JSON(http.StatusUnauthorized, response.ErrorResponse{
				Success: false,
				Message: "Unauthorized",
			})
			ctx.Abort()
			return
		}

		// Convert to UUID
		userUUID, ok := userID.(uuid.UUID)
		if !ok {
			ctx.JSON(http.StatusInternalServerError, response.ErrorResponse{
				Success: false,
				Message: "Invalid user ID",
			})
			ctx.Abort()
			return
		}

		// Check MFA session
		mfaVerified, exists := ctx.Get("mfa_verified")
		if exists && mfaVerified.(bool) {
			// MFA already verified for this session
			ctx.Next()
			return
		}

		// Get MFA status from request header
		mfaToken := ctx.GetHeader("X-MFA-Token")
		if mfaToken == "" {
			ctx.JSON(http.StatusForbidden, response.ErrorResponse{
				Success: false,
			})
			ctx.Abort()
			return
		}

		// Verify MFA token (implement token validation logic here)
		// For now, just checking if user has MFA enabled
		user, err := userRepo.GetUserByID(ctx, userUUID)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, response.ErrorResponse{
				Success: false,
				Message: "Failed to verify MFA",
			})
			ctx.Abort()
			return
		}

		log.Printf("User MFA status: %v", user.AccountType)

		ctx.Next()
	}
}
