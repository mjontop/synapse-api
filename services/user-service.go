package services

import (
	"context"
	"fmt"
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

	_, isExistingUser := checkIsNewUser(ctx, userRepo, user, c)
	if isExistingUser {
		return
	}

	err := userRepo.CreateUser(ctx, user)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User created successfully", "user": user})

}

func checkIsNewUser(ctx context.Context, userRepo repositories.UserRepository, user models.User, c *gin.Context) (error, bool) {
	existingUserWithEmail, err := userRepo.GetUserByEmail(ctx, user.Email)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return nil, true
	}

	if existingUserWithEmail.Email == user.Email {
		error := fmt.Sprintf("User with email \"%s\" is already taken", user.Email)
		c.JSON(http.StatusBadRequest, gin.H{"error": error})
		return nil, true
	}

	existingUserWithUsername, err := userRepo.GetUserByUserName(ctx, user.Username)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return nil, true
	}

	if existingUserWithUsername.Username == user.Username {
		error := fmt.Sprintf("User with username \"%s\" is already taken", user.Username)
		c.JSON(http.StatusBadRequest, gin.H{"error": error})
		return nil, true
	}
	return err, false
}
