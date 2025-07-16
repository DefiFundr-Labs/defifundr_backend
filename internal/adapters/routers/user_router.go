package routers

import (
	"github.com/demola234/defifundr/internal/adapters/handlers"
	"github.com/gin-gonic/gin"
)

func RegisterUserRoutes(rg *gin.RouterGroup, handler *handlers.UserHandler, authMiddleware gin.HandlerFunc) {
	users := rg.Group("/users")
	users.Use(authMiddleware)

	{
		users.GET("/profile", handler.GetProfile)
		// users.PUT("/profile", handler.UpdateProfile)
		users.POST("/change-password", handler.ChangePassword)

	}
}
