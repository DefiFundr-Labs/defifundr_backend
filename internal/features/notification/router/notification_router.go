package notificationrouter

import (
	notificationhandler "github.com/demola234/defifundr/internal/features/notification/handler"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(rg *gin.RouterGroup, handler *notificationhandler.Handler, authMiddleware gin.HandlerFunc) {
	g := rg.Group("/notifications")
	g.Use(authMiddleware)
	{
		g.GET("", handler.NotImplemented)
		g.POST("", handler.NotImplemented)
		g.GET("/:id", handler.NotImplemented)
		g.PUT("/:id", handler.NotImplemented)
		g.DELETE("/:id", handler.NotImplemented)
	}
}
