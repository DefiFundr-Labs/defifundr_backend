package userrouter

import (
	userhandler "github.com/demola234/defifundr/internal/features/user/handler"
	"github.com/gin-gonic/gin"
)

// RegisterRoutes registers user endpoints onto the provided router group.
func RegisterRoutes(rg *gin.RouterGroup, handler *userhandler.Handler, authMiddleware gin.HandlerFunc) {
	users := rg.Group("/users")
	users.Use(authMiddleware)
	{
		users.GET("/profile", handler.GetProfile)
		users.PUT("/profile", handler.UpdateProfile)
		users.POST("/change-password", handler.ChangePassword)
	}
}
