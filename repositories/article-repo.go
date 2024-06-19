package repositories

import (
	"context"
	"errors"
	"fmt"

	"github.com/mitchellh/mapstructure"
	"github.com/mjontop/synapse-api/db"
	"github.com/mjontop/synapse-api/lib/responses"
	"github.com/mjontop/synapse-api/models"
	"github.com/mjontop/synapse-api/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ArticleRepository interface {
	CreateArticle(ctx context.Context, article models.Article) error
	// GetArticles(ctx context.Context, filter bson.D, page int, pageSize int) ([]models.Article, error)
	GetArticleBySlug(ctx context.Context, slug string) (responses.ArticleResponseType, error)
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

func (repo *articleRepository) GetArticleBySlug(ctx context.Context, slug string) (responses.ArticleResponseType, error) {
	var article map[string]interface{}
	var articleResponse responses.ArticleResponseType

	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.D{{Key: "slug", Value: slug}, {Key: "isDeleted", Value: false}}}},
		{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "users"},
			{Key: "localField", Value: "author"},
			{Key: "foreignField", Value: "_id"},
			{Key: "as", Value: "author"},
		}}},
		{{Key: "$unwind", Value: bson.D{{Key: "path", Value: "$author"}}}},
		{{Key: "$project", Value: bson.D{
			{Key: "title", Value: 1},
			{Key: "id", Value: 1},
			{Key: "slug", Value: 1},
			{Key: "content", Value: 1},
			{Key: "description", Value: 1},
			{Key: "body", Value: 1},
			{Key: "tagList", Value: 1},
			{Key: "createdAt", Value: 1},
			{Key: "updatedAt", Value: 1},
			{Key: "author.username", Value: 1},
			{Key: "author.bio", Value: 1},
			{Key: "author.image", Value: 1},
			{Key: "author.fullname", Value: 1},
			{Key: "author.email", Value: 1},
		}}},
	}

	cursor, err := repo.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return articleResponse, fmt.Errorf("failed to execute aggregation: %v", err)
	}
	defer cursor.Close(ctx)

	var results []bson.M
	if err = cursor.All(ctx, &results); err != nil {
		return articleResponse, fmt.Errorf("failed to decode results: %v", err)
	}

	if len(results) == 0 {
		return articleResponse, errors.New("no article found with the given slug")
	}

	err = mapstructure.Decode(results[0], &article)
	if err != nil {
		return articleResponse, errors.New("no article found with the given slug")

	}

	articleResponse, err = convertToArticleResponse(article)
	if err != nil {
		return articleResponse, errors.New("no article found with the given slug")
	}
	return articleResponse, nil
}

func (repo *articleRepository) UpdateArticleByID(ctx context.Context, articleID primitive.ObjectID, update bson.D) error {
	update = bson.D{{Key: "$set", Value: update}}
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
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "isDeleted", Value: true}}}} // Setting isDeleted to true for soft delete
	_, err := repo.collection.UpdateByID(ctx, articleID, update)
	if err != nil {
		return err
	}
	return nil
}

func convertToArticleResponse(data map[string]interface{}) (responses.ArticleResponseType, error) {
	tagList := []string{}
	v := data["tagList"].(primitive.A)
	for _, tag := range v {
		tagList = append(tagList, tag.(string))
	}

	authorMap := data["author"].(primitive.M)

	createdAt, err := utils.ParseTime(data["createdAt"])
	if err != nil {
		return responses.ArticleResponseType{}, err
	}

	updatedAt, err := utils.ParseTime(data["updatedAt"])
	if err != nil {
		return responses.ArticleResponseType{}, err
	}

	return responses.ArticleResponseType{
		Title:               data["title"].(string),
		Slug:                data["slug"].(string),
		Description:         data["description"].(string),
		Body:                data["body"].(string),
		TagList:             tagList,
		PostCreationTimeUtc: createdAt,
		LastUpdatedAtUtc:    updatedAt,
		User: responses.UserDto{
			Email:    authorMap["email"].(string),
			Username: authorMap["username"].(string),
			Bio:      authorMap["bio"].(string),
			Image:    authorMap["image"].(string),
		},
	}, nil
}
