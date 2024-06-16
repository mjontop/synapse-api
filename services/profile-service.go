package services

import (
	"context"

	"github.com/mjontop/synapse-api/models"
	"github.com/mjontop/synapse-api/repositories"
)

func GetUserByUserName(username string) (models.User, error) {
	ctx := context.Background()
	userRepo := repositories.NewUserRepo()
	user, err := userRepo.GetUserByUserName(ctx, username)
	return user, err
}
