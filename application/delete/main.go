package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"os"
)

// LambdaEvent Lambda params
type LambdaEvent struct {
	ID string `json:"id"`
}

// LambdaResponse Lambda response model
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

// HandleLambdaEvent Lambda Handler
func HandleLambdaEvent(ctx context.Context, event LambdaEvent) (*LambdaResponse, error) {
	if event.ID != "" {
		bookID := event.ID

		db := dynamodb.New(getSession())

		query := &dynamodb.DeleteItemInput{
			TableName: aws.String(os.Getenv("LIBRARY_TABLE")),
			Key: map[string]*dynamodb.AttributeValue{
				"book_id": {
					S: aws.String(bookID),
				},
			},
		}

		_, err := db.DeleteItem(query)
		if err != nil {
			return nil, err
		}

		return &LambdaResponse{fmt.Sprintf("Book %s deleted", bookID)}, nil
	}

	return nil, errors.New("id is required")
}

func main() {
	lambda.Start(HandleLambdaEvent)
}