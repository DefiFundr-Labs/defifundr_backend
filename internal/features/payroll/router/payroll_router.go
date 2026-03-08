package payrollrouter

import (
	payrollhandler "github.com/demola234/defifundr/internal/features/payroll/handler"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(rg *gin.RouterGroup, handler *payrollhandler.Handler, authMiddleware gin.HandlerFunc) {
	g := rg.Group("/payroll")
	g.Use(authMiddleware)
	{
		g.GET("", handler.NotImplemented)
		g.POST("", handler.NotImplemented)
		g.GET("/:id", handler.NotImplemented)
		g.PUT("/:id", handler.NotImplemented)
		g.DELETE("/:id", handler.NotImplemented)
	}
}
