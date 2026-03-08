package adminrouter

import (
	adminhandler "github.com/demola234/defifundr/internal/features/admin/handler"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(rg *gin.RouterGroup, handler *adminhandler.Handler, authMiddleware gin.HandlerFunc) {
	g := rg.Group("/admin")
	g.Use(authMiddleware)
	{
		g.GET("", handler.NotImplemented)
		g.POST("", handler.NotImplemented)
		g.GET("/:id", handler.NotImplemented)
		g.PUT("/:id", handler.NotImplemented)
		g.DELETE("/:id", handler.NotImplemented)
	}
}
