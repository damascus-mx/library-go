package main

import (
	"github.com/damascus-mx/library-go/services/book/internal/shared/infrastructure/dependency"
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	bookHandler, close, err := dependency.InitializeBookHandler()
	defer close()
	if err != nil {
		panic(err)
	}

	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.POST("/v1/book", bookHandler.Save)

	server := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	err = server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}
