package services

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/mjontop/synapse-api/lib/requests"
	"github.com/mjontop/synapse-api/models"
	"github.com/mjontop/synapse-api/repositories"
	"github.com/mjontop/synapse-api/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
)

func HandleAddComment(c *gin.Context) {
	commentRepo := repositories.NewCommentRepository()
	slug := c.Param("slug")
	ctx := context.Background()

	user := c.MustGet("user").(models.User)

	var commentDto requests.CreateCommentRequestDto

	if err := c.BindJSON(&commentDto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": utils.ErrCommentCreate})
		return
	}

	articleRepo := repositories.NewArticleRepo()

	article, err := articleRepo.GetArticleBySlug(ctx, slug)

	if err != nil || article.Slug == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "article not found"})
		return
	}

	comment := models.NewComment(commentDto.Comment.Body, user.ID, slug)

	_, err = commentRepo.Create(ctx, comment)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": utils.ErrCommentCreate})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"comment": comment})
}

func HandleDeleteComment(c *gin.Context) {
	commentRepo := repositories.NewCommentRepository()
	commentId := c.Param("commentId")

	notFoundErr := gin.H{"error": "comment not found"}

	commentObjId, err := primitive.ObjectIDFromHex(commentId)
	if err != nil {
		c.JSON(http.StatusBadRequest, notFoundErr)
		return
	}

	ctx := context.Background()

	user := c.MustGet("user").(models.User)

	comment, err := commentRepo.GetByID(ctx, commentObjId)

	if err != nil {
		c.JSON(http.StatusInternalServerError, notFoundErr)
		return
	}

	if comment.Author != user.ID {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}

	err = commentRepo.Delete(ctx, commentObjId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": utils.ErrCommentCreate})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{"message": "Comment deleted successfully"})
}
