package services

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mjontop/synapse-api/lib/requests"
	"github.com/mjontop/synapse-api/lib/responses"
	"github.com/mjontop/synapse-api/models"
	"github.com/mjontop/synapse-api/repositories"
	"github.com/mjontop/synapse-api/utils"
)

func Register(c *gin.Context) {
	var userRequest requests.UserRegisterRequestType
	var user models.User

	ctx := context.Background()

	if err := c.ShouldBindJSON(&userRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user = userRequest.User

	userRepo := repositories.NewUserRepo()

	_, isExistingUser := checkIsNewUser(ctx, userRepo, user, c)
	if isExistingUser {
		return
	}

	hashedPassword, err := utils.HashPassword(user.Password)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	user.Password = hashedPassword

	err = userRepo.CreateUser(ctx, user)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user.Password = ""

	c.JSON(http.StatusOK, gin.H{"message": "User created successfully", "user": user})
}

func Login(c *gin.Context) {
	var loginUser requests.UserLoginRequestType
	var user models.User

	ctx := context.Background()

	if err := c.ShouldBindJSON(&loginUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	useRepo := repositories.NewUserRepo()

	if loginUser.User.Email != "" {
		existingUser, err := useRepo.GetUserByEmail(ctx, loginUser.User.Email)
		if err != nil || existingUser.Email == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Email"})
			return
		}
		user = existingUser

	} else if loginUser.User.Username != "" {
		existingUser, err := useRepo.GetUserByUserName(ctx, loginUser.User.Username)
		if err != nil || existingUser.Username == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Username"})
			return
		}
		user = existingUser

	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username or Email is Required"})
		return
	}

	isValidPassword := utils.CheckPassword(loginUser.User.Password, user.Password)

	if !isValidPassword {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Credentials"})
		return
	}

	token, err := utils.GenerateToken(user.Username)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error "})
		return
	}

	newLoginuser := responses.NewLoggedInUserResponse(user.Email, token, user.Username, user.Bio, user.Image)

	c.JSON(http.StatusOK, newLoginuser)
}

func RefreshCurrentLoggedUser(c *gin.Context) {
	currentUser := c.MustGet("user").(models.User)
	token, err := utils.GenerateToken(currentUser.Username)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error "})
		return
	}

	currentUserResponse := responses.NewLoggedInUserResponse(currentUser.Email, token, currentUser.Username, currentUser.Bio, currentUser.Image)

	c.JSON(http.StatusOK, currentUserResponse)
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
