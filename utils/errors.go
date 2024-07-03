package utils

import "errors"

var ErrArticleNotFound = errors.New("article not found")

var ErrCommentCreate = errors.New("bad request type for comment")
