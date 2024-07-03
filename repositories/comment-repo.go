package repositories

import (
	"context"
	"github.com/mjontop/synapse-api/db"
	"github.com/mjontop/synapse-api/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type CommentRepository interface {
	Create(ctx context.Context, comment *models.Comment) (*mongo.InsertOneResult, error)
	GetByID(ctx context.Context, id primitive.ObjectID) (*models.Comment, error)
	GetAllByArticleSlug(ctx context.Context, articleSlug string) ([]*models.Comment, error)
	Update(ctx context.Context, id primitive.ObjectID, update bson.M) (*mongo.UpdateResult, error)
	Delete(ctx context.Context, id primitive.ObjectID) (*mongo.DeleteResult, error)
}

type commentRepository struct {
	collection *mongo.Collection
}

func NewCommentRepository() CommentRepository {
	return &commentRepository{
		collection: db.GetCollection("comment"),
	}
}

func (r *commentRepository) Create(ctx context.Context, comment *models.Comment) (*mongo.InsertOneResult, error) {
	comment.ID = primitive.NewObjectID()
	comment.CreatedAt = time.Now()
	comment.UpdatedAt = time.Now()
	return r.collection.InsertOne(ctx, comment)
}

func (r *commentRepository) GetByID(ctx context.Context, id primitive.ObjectID) (*models.Comment, error) {
	var comment models.Comment
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&comment)
	return &comment, err
}

func (r *commentRepository) GetAllByArticleSlug(ctx context.Context, articleSlug string) ([]*models.Comment, error) {
	filter := bson.M{"articleSlug": articleSlug}
	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var comments []*models.Comment
	for cursor.Next(ctx) {
		var comment models.Comment
		if err := cursor.Decode(&comment); err != nil {
			return nil, err
		}
		comments = append(comments, &comment)
	}
	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return comments, nil
}

func (r *commentRepository) Update(ctx context.Context, id primitive.ObjectID, update bson.M) (*mongo.UpdateResult, error) {
	update["updatedAt"] = time.Now()
	return r.collection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": update})
}

func (r *commentRepository) Delete(ctx context.Context, id primitive.ObjectID) (*mongo.DeleteResult, error) {
	return r.collection.DeleteOne(ctx, bson.M{"_id": id})
}
