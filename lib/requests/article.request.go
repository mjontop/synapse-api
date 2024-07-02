package requests

type ArticleDto struct {
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Body        string   `json:"body"`
	TagList     []string `json:"tagList"`
}

type CreateArticleRequestDto struct {
	Article ArticleDto `json:"article"`
}

type UpdateArticleRequestDto struct {
	Article ArticleDto `json:"article"`
}

type GetArticleDto = CreateArticleRequestDto
