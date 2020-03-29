package model

import (
	"errors"
	"github.com/google/uuid"
	"log"
	"time"
)

// BookModel Library item entity
type BookModel struct {
	ID          string    `json:"book_id" docstore:"book_id"`
	Name        string    `json:"name" docstore:"name"`
	S3Url       *string   `json:"s3_url" docstore:"s3_url"`
	Authors     []*string `json:"authors" docstore:"authors"`
	Categories  []*string `json:"categories" docstore:"categories"`
	PublishedAt int64     `json:"published_at" docstore:"published_at"`
	CreatedAt   int64     `json:"created_at" docstore:"created_at"`
}

// NewBook Create a new book
func NewBook(name string, s3url *string, authors, categories []*string, publishedAt int64) *BookModel {
	return &BookModel{
		ID:          uuid.New().String(),
		Name:        name,
		S3Url:       s3url,
		Authors:     authors,
		Categories:  categories,
		PublishedAt: publishedAt,
		CreatedAt:   time.Now().Unix(),
	}
}

// TODO: Fix validations, emptiness
func (b *BookModel) Validate() error {
	log.Printf("%v", b)
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

	return nil
}
