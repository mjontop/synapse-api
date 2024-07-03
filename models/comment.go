package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Comment struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	CreatedAt   time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt   time.Time          `bson:"updatedAt" json:"updatedAt"`
	Body        string             `bson:"body" json:"body"`
	Author      primitive.ObjectID `bson:"author" json:"author"`
	ArticleSlug string             `bson:"articleSlug" json:"articleSlug"`
}

func NewComment(body string, author primitive.ObjectID, articleSlug string) *Comment {
	now := time.Now().UTC()
	return &Comment{
		CreatedAt:   now,
		UpdatedAt:   now,
		Body:        body,
		Author:      author,
		ArticleSlug: articleSlug,
	}
}
