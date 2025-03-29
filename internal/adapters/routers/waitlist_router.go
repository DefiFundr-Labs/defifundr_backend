package routers

import (
	"github.com/demola234/defifundr/internal/adapters/handlers"
	"github.com/gin-gonic/gin"
)

func RegisterWaitlistRoutes(rg *gin.RouterGroup, waitlistHandler *handlers.WaitlistHandler, authMiddleware gin.HandlerFunc) {
	// Public waitlist routes - accessible without authentication
	rg.POST("/waitlist", waitlistHandler.JoinWaitlist)
	
	// Admin routes - require authentication
	adminRoutes := rg.Group("/admin/waitlist")
	adminRoutes.Use(authMiddleware)
	{
		adminRoutes.GET("", waitlistHandler.ListWaitlist)
		adminRoutes.GET("/stats", waitlistHandler.GetWaitlistStats)
		adminRoutes.GET("/export", waitlistHandler.ExportWaitlist)
	}
}