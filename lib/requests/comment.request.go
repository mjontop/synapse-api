package requests

import "go.mongodb.org/mongo-driver/bson/primitive"

type CommentDto struct {
	Body        string             `json:"body"`
	AuthorID    primitive.ObjectID `json:"authorId"`
	ArticleSlug string             `json:"articleSlug"`
}

type CreateCommentRequestDto struct {
	Comment CommentDto `json:"comment"`
}

type UpdateCommentRequestDto struct {
	Comment CommentDto `json:"comment"`
}

type GetCommentDto = CreateCommentRequestDto
