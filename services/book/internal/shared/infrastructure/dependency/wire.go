// +build wireinject

package dependency

import (
	"context"
	"github.com/damascus-mx/library-go/services/book/internal/library/book/application"
	doc2 "github.com/damascus-mx/library-go/services/book/internal/library/book/domain/repository/doc"
	"github.com/damascus-mx/library-go/services/book/internal/library/book/infrastructure/persistence/doc"
	"github.com/damascus-mx/library-go/services/book/internal/presentation/delivery/http/handler"
	"github.com/damascus-mx/library-go/services/book/internal/shared/domain"
	"github.com/damascus-mx/library-go/services/book/internal/shared/infrastructure/persistence/service"
	"github.com/google/wire"
)

var dynamoSet = wire.NewSet(ProvideContext, domain.NewConfiguration, service.NewDynamoBookCollection)
var bookCaseSet = wire.NewSet(dynamoSet,
	doc.NewBookRepository,
	wire.Bind(new(doc2.IBookRepository), new(*doc.BookRepository)),
	application.NewBookUseCase,
)

// ProvideContext Get main context
func ProvideContext() context.Context {
	return context.Background()
}

func InitializeBookHandler() (*handler.BookHandler, func(), error) {
	wire.Build(wire.NewSet(bookCaseSet,
		handler.NewBookHandler,
	))

	return &handler.BookHandler{}, nil, nil
}
