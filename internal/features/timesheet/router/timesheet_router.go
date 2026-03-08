package timesheetrouter

import (
	timesheethandler "github.com/demola234/defifundr/internal/features/timesheet/handler"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(rg *gin.RouterGroup, handler *timesheethandler.Handler, authMiddleware gin.HandlerFunc) {
	g := rg.Group("/timesheets")
	g.Use(authMiddleware)
	{
		g.GET("", handler.NotImplemented)
		g.POST("", handler.NotImplemented)
		g.GET("/:id", handler.NotImplemented)
		g.PUT("/:id", handler.NotImplemented)
		g.DELETE("/:id", handler.NotImplemented)
	}
}
