package service

import (
	"context"
	"fmt"
	"github.com/damascus-mx/library-go/services/book/internal/shared/domain"
	"gocloud.dev/docstore"
	_ "gocloud.dev/docstore/awsdynamodb"
	"log"
)

func NewDynamoBookCollection(ctx context.Context, config *domain.Configuration) (*docstore.Collection, func(), error) {
	collection, err := docstore.OpenCollection(ctx, fmt.Sprintf("dynamodb://%s?partition_key=book_id&allow_scans=true", config.TableName))
	if err != nil {
		return nil, nil, err
	}

	cleanup := func() {
		if err := collection.Close(); err != nil {
			log.Print(err)
		}
	}
	return collection, cleanup, nil
}
