package repositories

import (
	"context"
	"errors"

	"github.com/mjontop/synapse-api/db"
	"github.com/mjontop/synapse-api/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user models.User) error
	GetUserByEmail(ctx context.Context, email string) (models.User, error)
	GetUserByUserName(ctx context.Context, username string) (models.User, error)
	UpdateUserById(user models.User) error
	GetUserById(ctx context.Context, id primitive.ObjectID) (models.User, error)
}

type userRepository struct {
	collection *mongo.Collection
}

func NewUserRepo() UserRepository {
	return &userRepository{
		collection: db.GetCollection("users"),
	}
}

func (repo *userRepository) GetUserByEmail(ctx context.Context, email string) (models.User, error) {
	var user models.User
	err := repo.collection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return user, nil
		}
		return user, err
	}
	return user, nil
}

func (repo *userRepository) GetUserById(ctx context.Context, id primitive.ObjectID) (models.User, error) {
	var user models.User
	err := repo.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return user, nil
		}
		return user, err
	}
	return user, nil
}

func (repo *userRepository) GetUserByUserName(ctx context.Context, username string) (models.User, error) {
	var user models.User
	err := repo.collection.FindOne(ctx, bson.M{"username": username}).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return user, nil
		}
		return user, err
	}
	return user, nil
}

func (repo *userRepository) CheckIfUserExists(ctx context.Context, user models.User) (bool, error) {
	emailUser, err := repo.GetUserByEmail(ctx, user.Email)
	if err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
		return false, err
	}
	if emailUser.Email == user.Email {
		return true, nil // Email Existed
	}

	usernameUser, err := repo.GetUserByUserName(ctx, user.Username)
	if err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
		return false, err
	}
	if usernameUser.Username == user.Username {
		return true, nil // Username exitsted
	}

	return false, nil
}

func (repo *userRepository) CreateUser(ctx context.Context, user models.User) error {

	userExists, err := repo.CheckIfUserExists(ctx, user)

	if err != nil {
		return err
	}

	if userExists {
		return errors.New("user already Exists")
	}

	user.ID = primitive.NewObjectID()

	_, err = repo.collection.InsertOne(ctx, user)
	return err
}

func (repo *userRepository) UpdateUserById(user models.User) error {

	ctx := context.Background()

	userExists, err := repo.CheckIfUserExists(ctx, user)

	if err != nil || !userExists {
		return err
	}

	dataToBeUpdated := bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "bio", Value: user.Bio},
			{Key: "image", Value: user.Image},
			{Key: "fullname", Value: user.FullName},
		}},
	}

	_, err = repo.collection.UpdateByID(ctx, user.ID, dataToBeUpdated)
	if err != nil {
		return err
	}

	return nil
}
