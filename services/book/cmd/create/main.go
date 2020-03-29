package main

import (
	"context"
	ginadapter "github.com/awslabs/aws-lambda-go-api-proxy/gin"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

var ginLambda *ginadapter.GinLambda

func init() {
	ginLambda = StartMux()
}

func main() {
	lambda.Start(Handler)
}

func Handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// If no name is provided in the HTTP request body, throw an error
	return ginLambda.ProxyWithContext(ctx, req)
}
