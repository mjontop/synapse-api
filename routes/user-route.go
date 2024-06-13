package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// SetupUserRoutes sets up routes related to users
func SetupUserRoutes(router *gin.RouterGroup) {
	userRoutes := router.Group("/users")
	{
		userRoutes.POST("/", func(c *gin.Context) {
			c.JSON(http.StatusCreated, gin.H{
				"message": "Register User route",
			})
		})

		userRoutes.POST("/login", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"message": "Login User route",
			})
		})

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
