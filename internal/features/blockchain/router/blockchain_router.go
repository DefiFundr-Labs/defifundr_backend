package blockchainrouter

import (
	blockchainhandler "github.com/demola234/defifundr/internal/features/blockchain/handler"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(rg *gin.RouterGroup, handler *blockchainhandler.Handler, authMiddleware gin.HandlerFunc) {
	g := rg.Group("/blockchain")
	g.Use(authMiddleware)
	{
		g.GET("", handler.NotImplemented)
		g.POST("", handler.NotImplemented)
		g.GET("/:id", handler.NotImplemented)
		g.PUT("/:id", handler.NotImplemented)
		g.DELETE("/:id", handler.NotImplemented)
	}
}
