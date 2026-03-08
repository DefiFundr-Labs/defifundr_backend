package waitlistrouter

import (
	waitlisthandler "github.com/demola234/defifundr/internal/features/waitlist/handler"
	"github.com/gin-gonic/gin"
)

// RegisterRoutes wires all waitlist endpoints onto the provided router group.
func RegisterRoutes(rg *gin.RouterGroup, handler *waitlisthandler.Handler, authMiddleware gin.HandlerFunc) {
	// Public: no authentication required
	rg.POST("/waitlist", handler.JoinWaitlist)

	// Admin: authentication + role check inside handler
	admin := rg.Group("/admin/waitlist")
	admin.Use(authMiddleware)
	{
		admin.GET("", handler.ListWaitlist)
		admin.GET("/stats", handler.GetWaitlistStats)
		admin.GET("/export", handler.ExportWaitlist)
	}
}
