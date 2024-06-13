package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/mjontop/synapse-api/db"
)

func main() {
	// Connect to DB
	db.ConnectDB()

	router := gin.Default()

	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Hello, Gin!",
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
