package walletrouter

import (
	wallethandler "github.com/demola234/defifundr/internal/features/wallet/handler"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(rg *gin.RouterGroup, handler *wallethandler.Handler, authMiddleware gin.HandlerFunc) {
	g := rg.Group("/wallets")
	g.Use(authMiddleware)
	{
		g.GET("", handler.NotImplemented)
		g.POST("", handler.NotImplemented)
		g.GET("/:id", handler.NotImplemented)
		g.PUT("/:id", handler.NotImplemented)
		g.DELETE("/:id", handler.NotImplemented)
	}
}
