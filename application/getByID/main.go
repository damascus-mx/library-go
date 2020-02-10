package main

import (
	"context"
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/damascus-mx/library-go/application/getByID/domain/model"
)

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

// HandleLambdaEvent AWS Lambda handler
func HandleLambdaEvent(ctx context.Context, req events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	db := dynamodb.New(getSession())

	idQuery := req.PathParameters["id"]
	if idQuery == "" {
		return proxyResponseBuilder("id is required", http.StatusBadRequest)
	}

	var queryInput = &dynamodb.QueryInput{
		TableName: aws.String(os.Getenv("LIBRARY_TABLE")),
		KeyConditions: map[string]*dynamodb.Condition{
			"book_id" : {
				ComparisonOperator: aws.String("EQ"),
				AttributeValueList: []*dynamodb.AttributeValue{
					{
						S: aws.String(idQuery),
					},
				},
			},
		},
	}

	result, err := db.Query(queryInput)
	if err != nil {
		return proxyResponseBuilder(err.Error(), http.StatusInternalServerError)
	}

	var books []model.BookModel
	err = dynamodbattribute.UnmarshalListOfMaps(result.Items, &books)
	if err != nil {
		return proxyResponseBuilder(err.Error(), http.StatusInternalServerError)
	}

	if len(books) < 1 {
		return proxyResponseBuilder("book not found", http.StatusNotFound)
	}

	bookJSON, _ := json.Marshal(&books[0])

	return &events.APIGatewayProxyResponse{
		StatusCode:        http.StatusOK,
		Headers:           map[string]string{
			"Access-Control-Allow-Origin": "*",
		},
		MultiValueHeaders: nil,
		Body:             string(bookJSON),
		IsBase64Encoded:   false,
	}, nil
}

func main() {
	lambda.Start(HandleLambdaEvent)
}
