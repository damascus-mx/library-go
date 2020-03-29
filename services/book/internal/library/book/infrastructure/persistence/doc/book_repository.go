package doc

import (
	"context"
	"errors"
	"github.com/damascus-mx/library-go/services/book/internal/library/book/domain/entity"
	"gocloud.dev/docstore"
	"io"
	"time"
)

type BookRepository struct {
	ctx        context.Context
	collection *docstore.Collection
}

func NewBookRepository(ctx context.Context, collection *docstore.Collection) *BookRepository {
	return &BookRepository{ctx, collection}
}

func (b *BookRepository) Save(book *entity.BookEntity) error {
	ctx, cancel := context.WithTimeout(b.ctx, 30*time.Second)
	defer cancel()

	if existingBook := b.FetchByName(book.Name); existingBook != nil {
		return errors.New("book already exists")
	}

	return b.collection.Create(ctx, book)
}

func (b *BookRepository) FetchByName(name string) *entity.BookEntity {
	ctx, cancel := context.WithTimeout(b.ctx, 30*time.Second)
	defer cancel()

	iter := b.collection.Query().Where("name", "=", name).Limit(1).Get(ctx)
	defer iter.Stop()

	existingBook := new(entity.BookEntity)
	for {
		err := iter.Next(ctx, existingBook)
		if err == io.EOF {
			break
		} else if err != nil {
			return nil
		} else {
			return existingBook
		}
	}

	return nil
}
