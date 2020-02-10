package main

import (
	"context"
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

// BookModel Library item model
type BookModel struct {
	ID          string    `json:"book_id"`
	Name        string    `json:"name"`
	S3Url		*string	  `json:"s3_url"`
	Authors     []*string `json:"authors"`
	Categories  []*string `json:"categories"`
	PublishedAt int64     `json:"published_at"`
	CreatedAt   int64     `json:"created_at"`
}

func proxyResponseBuilder(messageStr string, status int) (*events.APIGatewayProxyResponse, error) {
	message := struct {
		Message string
	}{ messageStr }

	jsonMsg, _ := json.Marshal(message)

	return &events.APIGatewayProxyResponse{
		StatusCode:        status,
		Headers:           map[string]string{
			"Access-Control-Allow-Origin": "*",
		},
		MultiValueHeaders: nil,
		Body:             string(jsonMsg),
		IsBase64Encoded:   false,
	}, nil
}

func getSession() *session.Session {
	config := &aws.Config{
		Region: aws.String("us-east-1"),
	}

	sess := session.Must(session.NewSession(config))

	return sess
}

// HandleLambdaEvent Lambda handler
func HandleLambdaEvent(ctx context.Context, req events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	db := dynamodb.New(getSession())

	params := &dynamodb.ScanInput{
		TableName: aws.String(os.Getenv("LIBRARY_TABLE")),
	}

	result, err := db.Scan(params)
	if err != nil {
		return proxyResponseBuilder(err.Error(), http.StatusInternalServerError)
	}

	var books []BookModel
	err = dynamodbattribute.UnmarshalListOfMaps(result.Items, &books)
	if err != nil {
		return proxyResponseBuilder(err.Error(), http.StatusInternalServerError)
	} else if len(books) == 0 {
		return proxyResponseBuilder("books not found", http.StatusNotFound)
	}

	booksJSON, _ := json.Marshal(books)

	return proxyResponseBuilder(string(booksJSON), http.StatusOK)
}

func main() {
	lambda.Start(HandleLambdaEvent)
}
