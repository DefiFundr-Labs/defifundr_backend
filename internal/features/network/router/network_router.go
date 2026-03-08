package networkrouter

import (
	networkhandler "github.com/demola234/defifundr/internal/features/network/handler"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(rg *gin.RouterGroup, handler *networkhandler.Handler, authMiddleware gin.HandlerFunc) {
	g := rg.Group("/networks")
	g.Use(authMiddleware)
	{
		g.GET("", handler.NotImplemented)
		g.POST("", handler.NotImplemented)
		g.GET("/:id", handler.NotImplemented)
		g.PUT("/:id", handler.NotImplemented)
		g.DELETE("/:id", handler.NotImplemented)
	}
}
