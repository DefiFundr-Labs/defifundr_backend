package invoicerouter

import (
	invoicehandler "github.com/demola234/defifundr/internal/features/invoice/handler"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(rg *gin.RouterGroup, handler *invoicehandler.Handler, authMiddleware gin.HandlerFunc) {
	g := rg.Group("/invoices")
	g.Use(authMiddleware)
	{
		g.GET("", handler.NotImplemented)
		g.POST("", handler.NotImplemented)
		g.GET("/:id", handler.NotImplemented)
		g.PUT("/:id", handler.NotImplemented)
		g.DELETE("/:id", handler.NotImplemented)
	}
}
