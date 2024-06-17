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

type ArticleRepository interface {
	CreateArticle(ctx context.Context, article models.Article) error
	// GetArticles(ctx context.Context, filter bson.D, page int, pageSize int) ([]models.Article, error)
	GetArticleByID(ctx context.Context, articleID primitive.ObjectID) (models.Article, error)
	UpdateArticleByID(ctx context.Context, articleID primitive.ObjectID, update bson.D) error
	DeleteArticleByID(ctx context.Context, articleID primitive.ObjectID) error
}

type articleRepository struct {
	collection *mongo.Collection
}

func NewArticleRepo() ArticleRepository {
	return &articleRepository{
		collection: db.GetCollection("articles"),
	}
}

func (repo *articleRepository) CreateArticle(ctx context.Context, article models.Article) error {
	article.ID = primitive.NewObjectID()
	_, err := repo.collection.InsertOne(ctx, article)
	return err
}

// func (repo *articleRepository) GetArticles(ctx context.Context, filter bson.D, page int, pageSize int) ([]models.Article, error) {
// 	skip := (page - 1) * pageSize
// 	sort := bson.D{{Key: "createdAt", Value: -1}} // Sort by creation date descending
// 	pipeline := bson.A{
// 		{Key: "$match", Value: filter},
// 		{Key: "$sort", Value: sort},
// 		{Key: "$skip", Value: skip},
// 		{Key: "$limit", Value: pageSize},
// 	}
// 	cursor, err := repo.collection.Aggregate(ctx, pipeline)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer cursor.Close(ctx)

// 	var articles []models.Article
// 	for cursor.Next(ctx) {
// 		var article models.Article
// 		if err := cursor.Decode(&article); err != nil {
// 			return nil, err
// 		}
// 		articles = append(articles, article)
// 	}
// 	return articles, nil
// }

func (repo *articleRepository) GetArticleByID(ctx context.Context, articleID primitive.ObjectID) (models.Article, error) {
	var article models.Article
	filter := bson.D{{Key: "_id", Value: articleID}, {Key: "isDeleted", Value: false}} // Filter for non-deleted articles
	err := repo.collection.FindOne(ctx, filter).Decode(&article)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return article, errors.New("article not found")
		}
		return article, err
	}
	return article, nil
}

func (repo *articleRepository) UpdateArticleByID(ctx context.Context, articleID primitive.ObjectID, update bson.D) error {
	update = bson.D{{Key: "$set", Value: update}} // Wrap update data in $set
	result, err := repo.collection.UpdateByID(ctx, articleID, update)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return errors.New("article not found")
	}
	return nil
}

func (repo *articleRepository) DeleteArticleByID(ctx context.Context, articleID primitive.ObjectID) error {
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "isDeleted", Value: true}}}} // Set isDeleted to true for soft delete
	_, err := repo.collection.UpdateByID(ctx, articleID, update)
	if err != nil {
		return err
	}
	return nil
}
