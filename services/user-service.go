package services

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mjontop/synapse-api/models"
	"github.com/mjontop/synapse-api/repositories"
)

func CreateUser(c *gin.Context) {
	var user models.User

	ctx := context.Background()

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userRepo := repositories.NewUserRepo()

	err := userRepo.CreateUser(ctx, user)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User created successfully", "user": user})

	// userRepo := db.GetCollection("users")

	// var existingUser models.User

	// err := userRepo.FindOne(ctx, bson.M{"email": user.Email}).Decode(&existingUser)

	// if err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
	// 	return
	// }

	// if existingUser.Email == user.Email {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("User with email \"%s\" is already registered", user.Email)})
	// 	return
	// }

	// _, err = userRepo.InsertOne(ctx, user)
	// if err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": err})
	// 	return
	// }

	// c.JSON(http.StatusOK, gin.H{"message": "User created successfully", "user": user})
}
