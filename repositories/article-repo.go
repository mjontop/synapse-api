package repositories

import (
	"context"
	"fmt"
	"github.com/mjontop/synapse-api/lib/requests"
	"time"

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
	GetPaginatedArticles(ctx context.Context) ([]responses.ArticleResponseType, error)
	GetArticleBySlug(ctx context.Context, slug string) (responses.ArticleResponseType, error)
	GetUsersArticleBySlug(ctx context.Context, slug string, currentUserId primitive.ObjectID) (models.Article, error)
	UpdateArticleBySlug(ctx context.Context, slug string, articleToBeUpdate *requests.ArticleDto, userId primitive.ObjectID) error
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

func (repo *articleRepository) GetPaginatedArticles(ctx context.Context) ([]responses.ArticleResponseType, error) {
	var articlesResponse []responses.ArticleResponseType
	var articles []map[string]interface{}

	pipeline := createArticlePipeline(bson.D{})

	cursor, err := repo.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return articlesResponse, fmt.Errorf("failed to execute aggregation: %v", err)
	}
	defer cursor.Close(ctx)

	var results []bson.M
	if err = cursor.All(ctx, &results); err != nil {
		return articlesResponse, fmt.Errorf("failed to decode results: %v", err)
	}

	err = mapstructure.Decode(results, &articles)
	if err != nil {
		return articlesResponse, utils.ErrArticleNotFound
	}

	if len(results) == 0 {
		return articlesResponse, utils.ErrArticleNotFound
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	for _, article := range articles {
		articleResponse, err := convertToArticleResponse(article)
		if err != nil {
			return articlesResponse, utils.ErrArticleNotFound
		}
		articlesResponse = append(articlesResponse, articleResponse)
	}

	return articlesResponse, nil
}

func (repo *articleRepository) GetArticleBySlug(ctx context.Context, slug string) (responses.ArticleResponseType, error) {
	var article map[string]interface{}
	var articleResponse responses.ArticleResponseType

	slugFilter := bson.D{{Key: "slug", Value: slug}}

	pipeline := createArticlePipeline(slugFilter)

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
		return articleResponse, utils.ErrArticleNotFound
	}

	err = mapstructure.Decode(results[0], &article)
	if err != nil {
		return articleResponse, utils.ErrArticleNotFound
	}

	articleResponse, err = convertToArticleResponse(article)
	if err != nil {
		return articleResponse, utils.ErrArticleNotFound
	}
	return articleResponse, nil
}

func (repo *articleRepository) DeleteArticleByID(ctx context.Context, articleID primitive.ObjectID) error {
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "isDeleted", Value: true}}}}
	_, err := repo.collection.UpdateByID(ctx, articleID, update)
	if err != nil {
		return err
	}
	return nil
}

func (repo *articleRepository) GetUsersArticleBySlug(ctx context.Context, slug string, currentUserId primitive.ObjectID) (models.Article, error) {
	var article models.Article

	slugFilter := bson.D{
		{Key: "slug", Value: slug},
		{Key: "author", Value: currentUserId},
		{Key: "isDeleted", Value: false},
	}

	err := repo.collection.FindOne(ctx, slugFilter).Decode(&article)

	if err != nil {
		return article, utils.ErrArticleNotFound
	}

	return article, nil
}

func (repo *articleRepository) UpdateArticleBySlug(ctx context.Context, slug string, articleToBeUpdate *requests.ArticleDto, userId primitive.ObjectID) error {
	_, err := repo.GetUsersArticleBySlug(ctx, slug, userId)
	if err != nil {
		return utils.ErrArticleNotFound
	}

	update := bson.M{}

	if articleToBeUpdate.Title != "" {
		update["title"] = articleToBeUpdate.Title
	}

	if articleToBeUpdate.Body != "" {
		update["body"] = articleToBeUpdate.Body
	}

	if articleToBeUpdate.Description != "" {
		update["description"] = articleToBeUpdate.Description
	}

	if cap(articleToBeUpdate.TagList) > 0 {
		update["tagList"] = articleToBeUpdate.TagList
	}

	update["updatedAt"] = time.Now().UTC()

	if len(update) == 0 {
		return nil
	}

	filter := bson.D{{"slug", slug}}

	_, err = repo.collection.UpdateOne(ctx, filter, bson.M{"$set": update})

	return err
}

func convertToArticleResponse(data map[string]interface{}) (responses.ArticleResponseType, error) {
	var tagList []string
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

func createArticlePipeline(filter bson.D) mongo.Pipeline {
	baseFilter := bson.D{{Key: "isDeleted", Value: false}}

	combinedFilter := append(filter, baseFilter...)

	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: combinedFilter}},
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

	return pipeline
}
