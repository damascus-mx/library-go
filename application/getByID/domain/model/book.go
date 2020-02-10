package model

// BookModel Library item model
type BookModel struct {
	ID          string    `json:"book_id"`
	Name        string    `json:"name"`
	S3Url		*string	  `json:"s3_url"`
	Authors     []*string `json:"authors"`
	Categories  []*string `json:"categories"`
	PublishedAt int64     `json:"published_at"`
	CreatedAt   int64     `json:"created_at"`
}
