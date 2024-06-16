package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/mjontop/synapse-api/services"
	"github.com/mjontop/synapse-api/utils"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization Failed"})
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		decodedToken, err := utils.DecodeToken(tokenString)

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Session Expired"})
			c.Abort()
			return
		}

		username := decodedToken["username"]

		user, err := services.GetUserByUserName(username)

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid User"})
			c.Abort()
			return
		}

		c.Set("user", user)

		c.Next()
	}
}
