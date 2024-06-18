package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/mjontop/synapse-api/middleware"
	"github.com/mjontop/synapse-api/services"
)

func SetupArticleRoutes(router *gin.RouterGroup) {
	articleRoutes := router.Group("/article")
	{
		articleRoutes.POST("/", middleware.AuthMiddleware(), services.HandleCreateArticle)
		articleRoutes.GET("/:slug", services.HandleGetArticleBySlug)

	}
}