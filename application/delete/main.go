package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"net/http"
	"os"
)

func getSession() *session.Session {
	config := &aws.Config{
		Region: aws.String("us-east-1"),
	}

	sess := session.Must(session.NewSession(config))

	return sess
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

// HandleLambdaEvent Lambda Handler
func HandleLambdaEvent(ctx context.Context, req events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	bookID := req.PathParameters["id"]
	if bookID != "" {
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
			return proxyResponseBuilder(err.Error(), http.StatusInternalServerError)
		}

		return proxyResponseBuilder(fmt.Sprintf("Book successfully %s deleted", bookID), http.StatusOK)
	}

	return proxyResponseBuilder("id is required", http.StatusBadRequest)
}

func main() {
	lambda.Start(HandleLambdaEvent)
}