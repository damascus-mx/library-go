
package model

import (
	"github.com/google/uuid"
	"time"
)

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

// NewBook Create a new book
func NewBook(name string, s3url *string, authors, categories []*string, publishedAt int64 ) *BookModel {
	return &BookModel{
		ID: uuid.New().String(),
		Name: name,
		S3Url: s3url,
		Authors: authors,
		Categories: categories,
		PublishedAt: publishedAt,
		CreatedAt: time.Now().Unix(),
	}
}