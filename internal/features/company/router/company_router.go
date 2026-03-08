package companyrouter

import (
	companyhandler "github.com/demola234/defifundr/internal/features/company/handler"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(rg *gin.RouterGroup, handler *companyhandler.Handler, authMiddleware gin.HandlerFunc) {
	g := rg.Group("/companies")
	g.Use(authMiddleware)
	{
		g.GET("", handler.NotImplemented)
		g.POST("", handler.NotImplemented)
		g.GET("/:id", handler.NotImplemented)
		g.PUT("/:id", handler.NotImplemented)
		g.DELETE("/:id", handler.NotImplemented)
	}
}
