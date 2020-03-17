package main

import (
	"context"
	"fmt"
	ginadapter "github.com/awslabs/aws-lambda-go-api-proxy/gin"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"gocloud.dev/docstore"
	_ "gocloud.dev/docstore/awsdynamodb"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"github.com/damascus-mx/library-go/services/book/application/create/domain/model"
	"github.com/gin-gonic/gin"
)

var ginLambda *ginadapter.GinLambda

func init() {
	log.Printf("Gin cold start")
	r := gin.Default()

	r.POST("/v1/book", SaveHandler)

	ginLambda = ginadapter.New(r)
}

func main() {
	lambda.Start(Handler)
}

func Handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// If no name is provided in the HTTP request body, throw an error
	return ginLambda.ProxyWithContext(ctx, req)
}

func SaveHandler(c *gin.Context) {
	bookJSON := new(model.BookRequest)

	err := c.BindJSON(bookJSON)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.Error{
			Err:  err,
			Type: http.StatusInternalServerError,
			Meta: err.Error(),
		})

		return
	}

	// Parse author form strings into idiomatic array
	authors := make([]*string, 0)
	for _, item := range bookJSON.Authors {
		authors = append(authors, &item)
	}

	// Parse categories form strings into idiomatic array
	categories := make([]*string,0)
	for _, item := range bookJSON.Categories {
		categories = append(categories, &item)
	}

	// Convert published_at to int64
	var publishedAt int64
	publishedAt = 0
	if bookJSON.PublishedAt != "" {
		convertedPublished, err := strconv.ParseInt(bookJSON.PublishedAt, 10, 64)
		if err != nil {
			publishedAt = 0
		}

		publishedAt = convertedPublished
	}

	book := model.NewBook(bookJSON.Name, nil, authors, categories, publishedAt)
	err = book.Validate()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.Error{
			Err:  err,
			Type: http.StatusBadRequest,
			Meta: err.Error(),
		})
		return
	}

	err = saveBook(book)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.Error{
			Err:  err,
			Type: http.StatusBadRequest,
			Meta: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("book %s created", book.ID),
	})
}

func saveBook(book *model.BookModel) error {
	ctxBackground := context.Background()
	ctx, cancel := context.WithTimeout(ctxBackground, 30*time.Second)
	defer cancel()

	coll, err := docstore.OpenCollection(ctx, fmt.Sprintf("dynamodb://%s?partition_key=book_id", os.Getenv("LIBRARY_TABLE")))
	if err != nil {
		return err
	}
	defer coll.Close()

	return coll.Create(ctx, book)
}
