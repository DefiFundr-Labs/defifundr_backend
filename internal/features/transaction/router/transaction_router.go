package transactionrouter

import (
	transactionhandler "github.com/demola234/defifundr/internal/features/transaction/handler"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(rg *gin.RouterGroup, handler *transactionhandler.Handler, authMiddleware gin.HandlerFunc) {
	g := rg.Group("/transactions")
	g.Use(authMiddleware)
	{
		g.GET("", handler.NotImplemented)
		g.POST("", handler.NotImplemented)
		g.GET("/:id", handler.NotImplemented)
		g.PUT("/:id", handler.NotImplemented)
		g.DELETE("/:id", handler.NotImplemented)
	}
}
