package model

// BookModel Library item model
type BookModel struct {
	Name        string    `json:"name"`
	Authors     []*string `json:"authors"`
	Categories  []*string `json:"categories"`
	PublishedAt int64     `json:"published_at"`
}

// NewBook Create a new book
func NewBook(name string, authors, categories []*string, publishedAt int64 ) *BookModel {
	return &BookModel{
		Name: name,
		Authors: authors,
		Categories: categories,
		PublishedAt: publishedAt,
	}
}