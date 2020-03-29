package main

import (
	ginadapter "github.com/awslabs/aws-lambda-go-api-proxy/gin"
	"github.com/damascus-mx/library-go/services/book/internal/shared/infrastructure/dependency"

	"github.com/gin-gonic/gin"
)

func StartMux() *ginadapter.GinLambda {
	// Resolve book Handler
	bookHandler, close, err := dependency.InitializeBookHandler()
	defer close()
	if err != nil {
		panic(err)
	}

	r := gin.Default()
	r.POST("/v1/book", bookHandler.Save)

	return ginadapter.New(r)
}
