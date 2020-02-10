package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/damascus-mx/library-go/application/create/domain/model"
	"os"
	"strconv"
)

// LambdaEvent AWS lambda params/events
type LambdaEvent struct{
	Name string `json:"name"`
	PublishedAt string `json:"published_at"`
	Authors []*string `json:"authors"`
	Categories []*string `json:"categories"`
}

// LambdaResponse AWS lambda response
type LambdaResponse struct {
	Message string `json:"message"`
}

func getSession() *session.Session {
	config := &aws.Config{
		Region: aws.String("us-east-1"),
	}

	sess := session.Must(session.NewSession(config))

	return sess
}


// HandleLambdaEvent AWS Lambda handler
func HandleLambdaEvent(ctx context.Context, event LambdaEvent) (*LambdaResponse, error) {
	if event.Name != "" {
		db := dynamodb.New(getSession())
		
		pubAt, err := strconv.ParseInt(event.PublishedAt, 10, 64)
		if err != nil {
			pubAt = 0
		}

		var authors []*string
		if len(event.Authors) == 0 {
			authors = nil
		} else {
			authors = event.Authors
		}

		var categories []*string
		if len(event.Categories) == 0 {
			categories = nil
		} else {
			categories = event.Categories
		}

		book := model.NewBook(event.Name, authors, categories, pubAt)

		bookMap, err := dynamodbattribute.MarshalMap(book)
		if err != nil {
			return nil, err
		}

		params := &dynamodb.PutItemInput{
			TableName: aws.String(os.Getenv("LIBRARY_TABLE")),
			Item: bookMap,
		}

		_, err = db.PutItem(params)
		if err != nil {
			return nil, err
		}

		return &LambdaResponse{fmt.Sprintf("User %s created", book.ID)}, nil
	}

	return nil, errors.New("name is required")
}

func main() {
	lambda.Start(HandleLambdaEvent)
}