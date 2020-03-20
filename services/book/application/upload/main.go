package main

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	ginadapter "github.com/awslabs/aws-lambda-go-api-proxy/gin"
	"github.com/damascus-mx/library-go/services/book/application/upload/domain/model"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gocloud.dev/blob"
	"gocloud.dev/docstore"
	_ "gocloud.dev/docstore/awsdynamodb"
	_ "gocloud.dev/blob/s3blob"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"strings"
	"time"
)

var ginLambda *ginadapter.GinLambda

func init() {
	log.Print("Gin cold start")
	mux := gin.Default()

	// Set MB limit for multipart form = 10 MB
	mux.MaxMultipartMemory = 10 << 20

	// Route mapping
	mux.POST("/v1/upload/book/:book_id", uploadBookHandler)

	ginLambda = ginadapter.New(mux)
}

func main() {
	lambda.Start(Handler)
}

func Handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// Proxy from default lambda handler to Gin engine router
	return ginLambda.ProxyWithContext(ctx, req)
}

func getCollectionConn() (*docstore.Collection, error) {
	ctxBackground := context.Background()
	ctx, cancel := context.WithTimeout(ctxBackground, 30*time.Second)
	defer cancel()

	coll, err := docstore.OpenCollection(ctx, fmt.Sprintf("dynamodb://%s?partition_key=book_id&allow_scans=true", os.Getenv("LIBRARY_TABLE")))
	if err != nil {
		return nil, err
	}

	return coll, nil
}

func uploadBookHandler(c *gin.Context) {
	// Get ID from params, verify if nil
	bookID := c.Param("book_id")
	if bookID == "" {
		err := errors.New("no book id received")
		c.JSON(http.StatusBadRequest, gin.Error{
			Err: err,
			Type: http.StatusBadRequest,
			Meta: err.Error(),
		})
		return
	}

	// GetBook, if nil we cannot start file uploading process
	book, err := getBook(bookID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.Error{
			Err: err,
			Type: http.StatusNotFound,
			Meta: err.Error(),
		})
		return
	}

	// Get file from multipart
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.Error{
			Err:  err,
			Type: http.StatusInternalServerError,
			Meta: err.Error(),
		})
		return
	}

	// Generate sanitized name
	fileName, err := sanitizeFileName(file.Filename)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.Error{
			Err:  err,
			Type: http.StatusBadRequest,
			Meta: err.Error(),
		})
		return
	}

	// Upload file to bucket
	err = uploadFile(fileName, file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.Error{
			Err:  err,
			Type: http.StatusInternalServerError,
			Meta: err.Error(),
		})
		return
	}

	// Update book
	err = updateBook(book, "https://cdn.damascus-engineering.com/damascus/ebooks/"+fileName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.Error{
			Err:  err,
			Type: http.StatusInternalServerError,
			Meta: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":fmt.Sprintf("file %s succesfuly uploaded", fileName),
	})
}

func getBook(bookID string) (*model.BookModel, error) {
	ctx := context.Background()
	coll, err := getCollectionConn()
	if err != nil {
		return nil, err
	}
	defer coll.Close()

	book := new(model.BookModel)
	book.ID = bookID
	err = coll.Get(ctx, book, "book_id")
	if err != nil {
		return nil, errors.New("book not found")
	}

	return book, nil
}

func getFileExtension(filename string) (string, error) {
	splatFileName := strings.Split(filename, ".")
	fileExtension := splatFileName[len(splatFileName) - 1]
	if len(splatFileName) < 2 {
		return "", errors.New("invalid file extension")
	} else if fileExtension != "pdf" {
		return "", errors.New("invalid file format, consider using pdf files")
	}

	return fileExtension, nil
}

func sanitizeFileName(filename string) (string, error) {
	// Generate file UUID for name uniqueness
	sanitizedName := uuid.New().String()

	// Get file extension
	fileExtension, err := getFileExtension(filename)
	if err != nil {
		return "", err
	}

	sanitizedName = fmt.Sprintf("%s.%s", sanitizedName, fileExtension)

	return sanitizedName, nil
}

func uploadFile(fileName string, fileHeader *multipart.FileHeader) error {
	ctx := context.Background()
	bucket, err := blob.OpenBucket(ctx, "s3://cdn.damascus-engineering.com?region=us-west-2")
	if err != nil {
		return err
	}
	bucket = blob.PrefixedBucket(bucket, "damascus/ebooks/")
	defer bucket.Close()

	// Write file to bucket
	// Create a cancelable context from the existing context.
	writeCtx, cancelWrite := context.WithCancel(ctx)
	defer cancelWrite()

	// Open file from file header
	file, err := fileHeader.Open()
	if err != nil {
		return err
	}
	defer file.Close()

	// Open bucket writer
	w, err := bucket.NewWriter(writeCtx, fileName, nil)
	if err != nil {
		return err
	}

	// Create new buffer and paste file data into it
	buf := bytes.NewBuffer(nil)
	_, err = io.Copy(buf, file)
	if err != nil {
		return err
	}

	// Write buffer's data
	_, writeErr := w.Write(buf.Bytes())
	// Close buffer writer securely
	if writeErr != nil {
		return writeErr
	}

	// Close bucket writer securely
	closeErr := w.Close()
	if closeErr != nil {
		return closeErr
	}

	return nil
}

func updateBook(book *model.BookModel, fileURL string) error {
	coll, err := getCollectionConn()
	if err != nil {
		return err
	}
	defer coll.Close()

	return coll.Actions().Update(book, docstore.Mods{"s3_url": fileURL}).Do(context.Background())
}