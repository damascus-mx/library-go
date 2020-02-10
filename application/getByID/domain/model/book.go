package model

// BookModel Library item model
type BookModel struct {
	ID          string    `json:"book_id"`
	Name        string    `json:"name"`
	Authors     []*string `json:"authors"`
	Categories  []*string `json:"categories"`
	PublishedAt int       `json:"published_at"`
	CreatedAt   int       `json:"created_at"`
}
