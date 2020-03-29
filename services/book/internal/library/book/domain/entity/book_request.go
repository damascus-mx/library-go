package entity

// BookRequest Book entity from request
type BookRequest struct {
	ID          string   `json:"book_id"`
	Name        string   `json:"name"`
	S3Url       string   `json:"s3_url"`
	Authors     []string `json:"authors"`
	Categories  []string `json:"categories"`
	PublishedAt string   `json:"published_at"`
	CreatedAt   string   `json:"created_at"`
}
