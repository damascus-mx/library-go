package main

import (
	"context"
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/damascus-mx/library-go/services/book/application/update/domain/model"
	"github.com/google/uuid"
	"net/http"
	"os"
	"strconv"
)

// EventBody Library item model
type EventBody struct {
	Name        string    `json:"name"`
	Authors     []*string `json:"authors"`
	Categories  []*string `json:"categories"`
	PublishedAt string     `json:"published_at"`
}

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

	jsonMsg, err := json.Marshal(message)
	if err != nil {
		return nil, err
	}

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

func LambdaEventHandler(ctx context.Context, req events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	bookID := req.PathParameters["id"]
	if bookID == "" {
		return proxyResponseBuilder("missing id", http.StatusBadRequest)
	}
	_, err := uuid.Parse(bookID)
	if err != nil {
		return proxyResponseBuilder("not a valid id", http.StatusBadRequest)
	}

	var payload EventBody
	err = json.Unmarshal([]byte(req.Body), &payload)
	if err != nil {
		return proxyResponseBuilder("malformed body", http.StatusBadRequest)
	}

	if payload.Name != "" {
		db := dynamodb.New(getSession())

		pubAt, err := strconv.ParseInt(payload.PublishedAt, 10, 64)
		if err != nil {
			pubAt = 0
		}

		var authors []*string
		if len(payload.Authors) == 0 {
			authors = nil
		} else {
			authors =
				payload.Authors
		}

		var categories []*string
		if len(payload.Categories) == 0 {
			categories = nil
		} else {
			categories = payload.Categories
		}

		book := model.NewBook(payload.Name, authors, categories, pubAt)


		authorsAV, err := dynamodbattribute.MarshalList(book.Authors)
		if err != nil {
			return proxyResponseBuilder("failed authors convert", http.StatusInternalServerError)
		}

		categoriesAV, err := dynamodbattribute.MarshalList(book.Categories)
		if err != nil {
			return proxyResponseBuilder("failed categories convert", http.StatusInternalServerError)
		}

		params := &dynamodb.UpdateItemInput{
			TableName: aws.String(os.Getenv("LIBRARY_TABLE")),
			Key: map[string]*dynamodb.AttributeValue{
				"book_id": {
					S: aws.String(bookID),
				},
			},
			ExpressionAttributeNames: map[string]*string{
				"#nm": aws.String("name"),
			},
			UpdateExpression: aws.String("set #nm=:n, authors=:a, categories=:c, published_at=:p"),
			ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
				":n": {S: aws.String(book.Name)},
				":a": {L: authorsAV},
				":c": {L: categoriesAV},
				":p": {N: aws.String(strconv.Itoa(int(book.PublishedAt)))},
			},
			ReturnValues: aws.String(dynamodb.ReturnValueAllNew),
		}

		_, err = db.UpdateItem(params)
		if err != nil {
			return proxyResponseBuilder(err.Error(), http.StatusInternalServerError)
		}

		return proxyResponseBuilder("book successfully updated", http.StatusOK)
	}


	return proxyResponseBuilder("missing required field", http.StatusBadRequest)
}

func main() {
	lambda.Start(LambdaEventHandler)
}