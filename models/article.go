package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Article represents a blog article
type Article struct {
	ID                  primitive.ObjectID `bson:"_id,omitempty" json:"-"`
	Slug                string             `bson:"slug" json:"slug"`
	Title               string             `bson:"title" json:"title"`
	Description         string             `bson:"description" json:"description"`
	Body                string             `bson:"body" json:"body"`
	TagList             []string           `bson:"tagList" json:"tagList"`
	AuthorID            primitive.ObjectID `bson:"author" json:"author"` // Reference the User ID
	PostCreationTimeUtc time.Time          `bson:"createdAt" json:"createdAt"`
	LastUpdatedAtUtc    time.Time          `bson:"updatedAt" json:"updatedAt"`
	IsDeleted           bool               `bson:"isDeleted" json:"-"`

	// For internal purpose only
	CreatedAtUtc time.Time `bson:"createdAt" json:"-"`
	UpdatedAtUtc time.Time `bson:"updatedAt" json:"-"`

	//Later will add follwings
	// Favorited      bool               `bson:"favorited" json:"favorited"`
	// FavoritesCount int                `bson:"favoritesCount" json:"favoritesCount"`
}
