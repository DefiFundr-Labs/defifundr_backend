package hrrouter

import (
	hrhandler "github.com/demola234/defifundr/internal/features/hr/handler"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(rg *gin.RouterGroup, handler *hrhandler.Handler, authMiddleware gin.HandlerFunc) {
	g := rg.Group("/hr")
	g.Use(authMiddleware)
	{
		g.GET("", handler.NotImplemented)
		g.POST("", handler.NotImplemented)
		g.GET("/:id", handler.NotImplemented)
		g.PUT("/:id", handler.NotImplemented)
		g.DELETE("/:id", handler.NotImplemented)
	}
}
