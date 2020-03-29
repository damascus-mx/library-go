package application

import (
	"github.com/damascus-mx/library-go/services/book/internal/library/book/domain/entity"
	"github.com/damascus-mx/library-go/services/book/internal/library/book/domain/repository/doc"
	"strconv"
)

type BookUseCase struct {
	bookRepository doc.IBookRepository
}

func NewBookUseCase(bookRepo doc.IBookRepository) *BookUseCase {
	return &BookUseCase{bookRepo}
}

func (b *BookUseCase) Create(bookRequest *entity.BookRequest) error {
	// Parse author form strings into idiomatic array
	authors := make([]*string, 0)
	for _, item := range bookRequest.Authors {
		authors = append(authors, &item)
	}

	// Parse categories form strings into idiomatic array
	categories := make([]*string, 0)
	for _, item := range bookRequest.Categories {
		categories = append(categories, &item)
	}

	// Convert published_at to int64
	var publishedAt int64
	publishedAt = 0
	if bookRequest.PublishedAt != "" {
		convertedPublished, err := strconv.ParseInt(bookRequest.PublishedAt, 10, 64)
		if err != nil {
			publishedAt = 0
		}

		publishedAt = convertedPublished
	}

	book := entity.NewBook(bookRequest.Name, nil, authors, categories, publishedAt)
	err := book.Validate()
	if err != nil {
		return err
	}

	return b.bookRepository.Save(book)
}
