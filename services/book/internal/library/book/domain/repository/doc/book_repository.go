package doc

import "github.com/damascus-mx/library-go/services/book/internal/library/book/domain/entity"

type IBookRepository interface {
	Save(book *entity.BookEntity) error
	FetchByName(name string) *entity.BookEntity
}
