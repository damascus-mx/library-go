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
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/damascus-mx/library-go/application/create/domain/model"
	"net/http"
	"os"
	"strconv"
)

func getSession() *session.Session {
	config := &aws.Config{
		Region: aws.String("us-east-1"),
	}

	sess := session.Must(session.NewSession(config))

	return sess
}

type EventBody struct {
	Name string `json:"name"`
	PublishedAt string `json:"published_at"`
	Authors []*string `json:"authors"`
	Categories []*string `json:"categories"`
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

// HandleLambdaEvent AWS Lambda handler
func HandleLambdaEvent(ctx context.Context, req events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	
	var body EventBody
	err := json.Unmarshal([]byte(req.Body), &body)
	if err != nil {
		return proxyResponseBuilder(err.Error(), http.StatusInternalServerError)
	}

	if body.Name != "" {
		db := dynamodb.New(getSession())
		
		pubAt, err := strconv.ParseInt(body.PublishedAt, 10, 64)
		if err != nil {
			pubAt = 0
		}

		var authors []*string
		if len(body.Authors) == 0 {
			authors = nil
		} else {
			authors = body.Authors
		}

		var categories []*string
		if len(body.Categories) == 0 {
			categories = nil
		} else {
			categories = body.Categories
		}

		book := model.NewBook(body.Name, nil, authors, categories, pubAt)

		bookMap, err := dynamodbattribute.MarshalMap(book)
		if err != nil {
			return proxyResponseBuilder(err.Error(), http.StatusInternalServerError)
		}

		params := &dynamodb.PutItemInput{
			TableName: aws.String(os.Getenv("LIBRARY_TABLE")),
			Item: bookMap,
		}

		_, err = db.PutItem(params)
		if err != nil {
			return proxyResponseBuilder(err.Error(), http.StatusInternalServerError)
		}

		message := struct {
			Message string
		}{ fmt.Sprintf("Book %s created", book.ID) }

		jsonMsg, _ := json.Marshal(message)

		return proxyResponseBuilder(string(jsonMsg), http.StatusOK)
	}

	return proxyResponseBuilder("name is required", http.StatusBadRequest)
}

func main() {
	lambda.Start(HandleLambdaEvent)
}