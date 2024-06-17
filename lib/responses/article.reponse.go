package responses

type ArticleResponseType struct {
	Title string `json:"title"`
	Slug  string `json:"slug"`

	Description string   `json:"description"`
	Body        string   `json:"body"`
	TagList     []string `json:"tagList"`
	Author      UserDto  `json:"user"`
}
