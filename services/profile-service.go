package services

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mjontop/synapse-api/lib/responses"
	"github.com/mjontop/synapse-api/models"
	"github.com/mjontop/synapse-api/repositories"
)

func GetProfileHandler(c *gin.Context) {
	username := c.Param("username")
	if username == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Username"})
		return
	}

	user, err := GetUserByUserName(username)

	userDto := responses.UserDto{
		Email:    user.Email,
		Token:    "",
		Username: user.Username,
		Bio:      user.Bio,
		Image:    user.Image,
	}

	profileResponse := responses.NewUsersProfileResponse(userDto)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusOK, profileResponse)
}

func GetUserByUserName(username string) (models.User, error) {
	ctx := context.Background()
	userRepo := repositories.NewUserRepo()
	user, err := userRepo.GetUserByUserName(ctx, username)
	return user, err
}
