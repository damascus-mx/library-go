package model

import (
	"github.com/google/uuid"
	"time"
)

// BookModel Library item model
type BookModel struct {
	ID          string    `json:"book_id"`
	Name        string    `json:"name"`
	Authors     []*string `json:"authors"`
	Categories  []*string `json:"categories"`
	PublishedAt int64      `json:"published_at"`
	CreatedAt   int64      `json:"created_at"`
}

// NewBook Create a new book
func NewBook(name string, authors, categories []*string, publishedAt int64 ) *BookModel {
	return &BookModel{
		ID: uuid.New().String(),
		Name: name,
		Authors: authors,
		Categories: categories,
		PublishedAt: publishedAt,
		CreatedAt: time.Now().Unix(),
	}
}