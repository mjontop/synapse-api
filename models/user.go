package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// User represents a user in the system
type User struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"-"`
	Username     string             `bson:"username" json:"username"`
	FullName     string             `bson:"fullname" json:"fullname"`
	Email        string             `bson:"email" json:"email"`
	Password     string             `bson:"password" json:"password,omitempty"`
	Bio          string             `bson:"bio" json:"bio"`
	Image        string             `bson:"image" json:"image"`
	CreatedAtUtc time.Time          `bson:"createdAt" json:"-"`
	UpdatedAtUtc time.Time          `bson:"updatedAt" json:"-"`
}

// NewUser creates a new instance of User
func NewUser(username, email, password, bio, image string) *User {
	return &User{
		Username: username,
		Email:    email,
		Password: password,
		Bio:      bio,
		Image:    image,
	}
}
