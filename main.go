package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mjontop/synapse-api/db"
	"github.com/mjontop/synapse-api/routes"
)

func main() {
	// Connect to DB
	db.ConnectDB()

	router := gin.Default()
	api := router.Group("/api")

	routes.SetupUserRoutes(api)
	routes.SetupArticleRoutes(api)

	router.GET("/healthcheck", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": time.Now(),
		})
	})

	address := os.Getenv("ADDRESS")

	err := router.Run(address)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	} else {
		println("Server is started at: ", address)
	}

}
