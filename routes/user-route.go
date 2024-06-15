package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mjontop/synapse-api/services"
)

func SetupUserRoutes(router *gin.RouterGroup) {
	userRoutes := router.Group("/users")
	{
		userRoutes.POST("/", services.Register)

		userRoutes.POST("/login", services.Login)

		userRoutes.GET("/current", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"message": "Get Current User route",
			})
		})

		userRoutes.PUT("/update", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"message": "Update User route",
			})
		})
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
