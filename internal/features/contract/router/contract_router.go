package contractrouter

import (
	contracthandler "github.com/demola234/defifundr/internal/features/contract/handler"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(rg *gin.RouterGroup, handler *contracthandler.Handler, authMiddleware gin.HandlerFunc) {
	g := rg.Group("/contracts")
	g.Use(authMiddleware)
	{
		g.GET("", handler.NotImplemented)
		g.POST("", handler.NotImplemented)
		g.GET("/:id", handler.NotImplemented)
		g.PUT("/:id", handler.NotImplemented)
		g.DELETE("/:id", handler.NotImplemented)
	}
}
