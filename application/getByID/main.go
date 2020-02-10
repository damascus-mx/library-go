package main

import (
	"context"
	"errors"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"os"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/damascus-mx/library-go/application/getByID/domain/model"
)

// LambdaEvent AWS lambda params/events
type LambdaEvent struct{
	ID string `json:"id"`
}

func getSession() *session.Session {
	config := &aws.Config{
		Region: aws.String("us-east-1"),
	}

	sess := session.Must(session.NewSession(config))

	return sess
}

// HandleLambdaEvent AWS Lambda handler
func HandleLambdaEvent(ctx context.Context, event LambdaEvent) (*model.BookModel, error) {
	db := dynamodb.New(getSession())

	idQuery := event.ID
	if idQuery == "" {
		return nil, errors.New("id not provided")
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
		return nil, err
	}

	var books []model.BookModel
	err = dynamodbattribute.UnmarshalListOfMaps(result.Items, &books)
	if err != nil {
		return nil, err
	}

	if len(books) < 1 {
		return nil, errors.New("book not found")
	}

	return &books[0], nil
}

func main() {
	lambda.Start(HandleLambdaEvent)
}
