package middleware

import (
	"net/http"
	"sync"
	"time"

	response "github.com/demola234/defifundr/internal/adapters/dto/response"
	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

// RateLimitMiddleware limits the number of authentication attempts
func RateLimitMiddleware(limit int, duration time.Duration) gin.HandlerFunc {
	// Create a limiter for each IP
	ipLimiters := make(map[string]*rate.Limiter)
	mu := &sync.Mutex{}

	return func(ctx *gin.Context) {
		ip := ctx.ClientIP()

		// Get the rate limiter for this IP
		mu.Lock()
		limiter, exists := ipLimiters[ip]
		if !exists {
			limiter = rate.NewLimiter(rate.Limit(limit)/rate.Limit(duration.Seconds()), limit)
			ipLimiters[ip] = limiter
		}
		mu.Unlock()

		// Check if this request is allowed
		if !limiter.Allow() {
			ctx.JSON(http.StatusTooManyRequests, response.ErrorResponse{
				Success: false,
				Message: "Rate limit exceeded. Please try again later.",
			})
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}
