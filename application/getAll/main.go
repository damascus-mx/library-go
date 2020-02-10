package main

import (
	"context"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

// BookModel Library item model
type BookModel struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Authors     []*string `json:"authors"`
	Categories  []*string `json:"categories"`
	PublishedAt int       `json:"published_at"`
	CreatedAt   int       `json:"created_at"`
}

// LambdaEvent Lambda Event (params)
type LambdaEvent struct {
	Name string `json:"name"`
	Age  string `json:"age"`
}

// LambdaResponse Lambda response
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

// HandleLambdaEvent Lambda handler
func HandleLambdaEvent(ctx context.Context, event LambdaEvent) ([]BookModel, error) {
	db := dynamodb.New(getSession())

	params := &dynamodb.ScanInput{
		TableName: aws.String(os.Getenv("LIBRARY_TABLE")),
	}

	result, err := db.Scan(params)
	if err != nil {
		return nil, err
	}

	var books []BookModel
	err = dynamodbattribute.UnmarshalListOfMaps(result.Items, &books)
	if err != nil {
		return nil, err
	}

	return books, nil
}

func main() {
	lambda.Start(HandleLambdaEvent)
}