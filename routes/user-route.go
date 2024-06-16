package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mjontop/synapse-api/middleware"
	"github.com/mjontop/synapse-api/services"
)

func SetupUserRoutes(router *gin.RouterGroup) {
	userRoutes := router.Group("/user")
	{
		userRoutes.GET("/", middleware.AuthMiddleware(), services.RefreshCurrentLoggedUser)
		userRoutes.PATCH("/", middleware.AuthMiddleware(), services.UpdateUser)

	}

	usersRoutes := router.Group("/users")
	{
		usersRoutes.POST("/", services.Register)

		usersRoutes.POST("/login", services.Login)
	}

	profileRoutes := router.Group("/profiles/:username")
	{
		profileRoutes.GET("/", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"message": "Get User Profile route",
			})
		})

		profileRoutes.POST("/follow", func(c *gin.Context) {
			c.JSON(http.StatusCreated, gin.H{
				"message": "Follow User route",
			})
		})

		profileRoutes.DELETE("/unfollow", func(c *gin.Context) {
			c.JSON(http.StatusCreated, gin.H{
				"message": "Unfollow User route",
			})
		})
	}
}
