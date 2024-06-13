package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Article represents a blog article
type Article struct {
	ID             primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Slug           string             `bson:"slug" json:"slug"`
	Title          string             `bson:"title" json:"title"`
	Description    string             `bson:"description" json:"description"`
	Body           string             `bson:"body" json:"body"`
	TagList        []string           `bson:"tagList" json:"tagList"`
	CreatedAt      time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt      time.Time          `bson:"updatedAt" json:"updatedAt"`
	Favorited      bool               `bson:"favorited" json:"favorited"`
	FavoritesCount int                `bson:"favoritesCount" json:"favoritesCount"`
	Author         Author             `bson:"author" json:"author"`
}

// Author represents the author of an article
type Author struct {
	Username  string `bson:"username" json:"username"`
	Bio       string `bson:"bio" json:"bio"`
	Image     string `bson:"image" json:"image"`
	Following bool   `bson:"following" json:"following"`
}
