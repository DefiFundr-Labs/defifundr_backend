package compliancerouter

import (
	compliancehandler "github.com/demola234/defifundr/internal/features/compliance/handler"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(rg *gin.RouterGroup, handler *compliancehandler.Handler, authMiddleware gin.HandlerFunc) {
	g := rg.Group("/compliance")
	g.Use(authMiddleware)
	{
		g.GET("", handler.NotImplemented)
		g.POST("", handler.NotImplemented)
		g.GET("/:id", handler.NotImplemented)
		g.PUT("/:id", handler.NotImplemented)
		g.DELETE("/:id", handler.NotImplemented)
	}
}
