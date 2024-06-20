package services

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mjontop/synapse-api/lib/requests"
	"github.com/mjontop/synapse-api/lib/responses"
	"github.com/mjontop/synapse-api/models"
	"github.com/mjontop/synapse-api/repositories"
	"github.com/mjontop/synapse-api/utils"
)

func HandleCreateArticle(c *gin.Context) {
	articleRepo := repositories.NewArticleRepo()
	var articleDto requests.CreateArticleRequestDto

	if err := c.BindJSON(&articleDto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := c.MustGet("user").(models.User)

	slug := utils.GenerateSlug(articleDto.Article.Title)

	article := models.Article{
		Title:               articleDto.Article.Title,
		Slug:                slug,
		Description:         articleDto.Article.Description,
		Body:                articleDto.Article.Body,
		TagList:             articleDto.Article.TagList,
		CreatedAtUtc:        time.Now().UTC(),
		LastUpdatedAtUtc:    time.Now().UTC(),
		PostCreationTimeUtc: time.Now().UTC(),
		UpdatedAtUtc:        time.Now().UTC(),
		AuthorID:            user.ID,
	}

	// Validate article data (optional)
	// You can add validation logic here using appropriate libraries
	// and return an error response if validation fails.

	ctx := context.Background()
	err := articleRepo.CreateArticle(ctx, article)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	articleResponse := responses.ArticleResponseType{
		Title:               articleDto.Article.Title,
		Slug:                slug,
		Description:         articleDto.Article.Description,
		Body:                articleDto.Article.Body,
		TagList:             articleDto.Article.TagList,
		PostCreationTimeUtc: article.PostCreationTimeUtc,
		LastUpdatedAtUtc:    article.LastUpdatedAtUtc,
		User: responses.UserDto{
			Email:    user.Email,
			Username: user.Username,
			Bio:      user.Bio,
			Image:    user.Image,
		},
	}

	c.JSON(http.StatusCreated, gin.H{"article": articleResponse})
}

func HandleGetArticleBySlug(c *gin.Context) {
	articleRepo := repositories.NewArticleRepo()

	slug := c.Param("slug")

	notFoundError := gin.H{"error": fmt.Sprintf("Article with slug \"%s\" does not exists", slug)}

	if slug == "" {
		c.JSON(http.StatusBadRequest, notFoundError)
	}

	ctx := context.Background()
	article, err := articleRepo.GetArticleBySlug(ctx, slug)
	if err != nil {

		if errors.Is(err, utils.ErrArticleNotFound) {
			c.JSON(http.StatusBadRequest, notFoundError)
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"article": article})
}

func HandleGetPaginatedArticles(c *gin.Context) {
	articleRepo := repositories.NewArticleRepo()
	ctx := context.Background()
	article, err := articleRepo.GetPaginatedArticles(ctx)

	if err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"articles": article})
}
