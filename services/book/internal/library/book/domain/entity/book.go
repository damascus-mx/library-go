package entity

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"time"
)

// BookModel Library item entity
type BookEntity struct {
	ID          string    `json:"book_id" docstore:"book_id"`
	Name        string    `json:"name" docstore:"name"`
	S3Url       *string   `json:"s3_url" docstore:"s3_url"`
	Authors     []*string `json:"authors" docstore:"authors"`
	Categories  []*string `json:"categories" docstore:"categories"`
	PublishedAt int64     `json:"published_at" docstore:"published_at"`
	CreatedAt   int64     `json:"created_at" docstore:"created_at"`
}

// NewBook Create a new book
func NewBook(name string, s3url *string, authors, categories []*string, publishedAt int64) *BookEntity {
	return &BookEntity{
		ID:          uuid.New().String(),
		Name:        name,
		S3Url:       s3url,
		Authors:     authors,
		Categories:  categories,
		PublishedAt: publishedAt,
		CreatedAt:   time.Now().Unix(),
	}
}

func (b *BookEntity) Validate() error {
	// name, author and category must be between 1 and 50
	// categories and authors limit's 5
	if b.ID == "" {
		return errors.New("invalid id")
	} else if b.Name == "" || len(b.Name) > 50 {
		return errors.New("invalid name")
	} else if b.Authors == nil || len(b.Authors) == 0 {
		return errors.New("at least one author is required")
	} else if b.Categories == nil || len(b.Categories) == 0 {
		return errors.New("at least one category is required")
	} else if len(b.Categories) > 5 {
		return errors.New("categories limit exceeded")
	} else if len(b.Authors) > 5 {
		return errors.New("authors limit exceeded")
	}

	for i, author := range b.Authors {
		if author == nil || len(*author) == 0 || len(*author) > 50 {
			return errors.New(fmt.Sprintf("author %d is not valid", i+1))
		}
	}

	for i, category := range b.Categories {
		if category == nil || len(*category) == 0 || len(*category) > 50 {
			return errors.New(fmt.Sprintf("category %d is not valid", i+1))
		}
	}

	return nil
}
