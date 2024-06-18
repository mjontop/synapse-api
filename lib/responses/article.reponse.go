package responses

import "time"

type ArticleResponseType struct {
	Title string `json:"title"`
	Slug  string `json:"slug"`

	Description         string    `json:"description"`
	Body                string    `json:"body"`
	TagList             []string  `json:"tagList"`
	PostCreationTimeUtc time.Time `bson:"createdAtUtc" json:"createdAt"`
	LastUpdatedAtUtc    time.Time `bson:"updatedAtUtc" json:"updatedAt"`
	User                UserDto   `json:"author"`
}
