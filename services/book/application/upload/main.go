package main

import (
	"context"
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/google/uuid"
	"io"
	"mime/multipart"
	"net/http"
	"strings"
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

func LambdaHandler(ctx context.Context, req events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {

	id := req.PathParameters["id"]
	if _, err := uuid.Parse(id); err != nil {
		return proxyResponseBuilder("invalid id", http.StatusBadRequest)
	}

	if req.Body == "" {
		return proxyResponseBuilder("missing required fields", http.StatusBadRequest)
	}

	/*
	_, err := base64.StdEncoding.DecodeString(req.Body)
	if err != nil {
		return proxyResponseBuilder("failed to decode file", http.StatusInternalServerError)
	}*/
	var maxFileSize int64
	maxFileSize = 20 * 1000 * 1000


	// fileType := req.Headers["Content-Type"]
	var r io.Reader

	mp := multipart.NewReader(r, req.Body)
	f, err := mp.ReadForm(maxFileSize)
	if err == io.EOF {
		return proxyResponseBuilder("failed to decode file", http.StatusInternalServerError)
	}
	x := f.Value["Content-Type"][0]

	return &events.APIGatewayProxyResponse{
		StatusCode:        200,
		Headers:           map[string]string{
			"Content-Type": strings.Split(req.Headers["Content-Type"], ";")[0],
		},
		MultiValueHeaders: nil,
		Body:              x,
		IsBase64Encoded:   false,
	}, nil



	// return proxyResponseBuilder(req.Body, http.StatusNotFound)
}

func main() {
	lambda.Start(LambdaHandler)
}