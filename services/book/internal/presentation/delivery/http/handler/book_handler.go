package handler

import (
	"errors"
	"github.com/damascus-mx/library-go/services/book/internal/library/book/application"
	"github.com/damascus-mx/library-go/services/book/internal/library/book/domain/entity"
	"github.com/gin-gonic/gin"
	"net/http"
)

type BookHandler struct {
	bookUseCase *application.BookUseCase
}

func NewBookHandler(bookUseCase *application.BookUseCase) *BookHandler {
	return &BookHandler{bookUseCase}
}

func (b *BookHandler) Save(c *gin.Context) {

	bookJSON := new(entity.BookRequest)
	err := c.BindJSON(bookJSON)
	if err != nil {
		err = errors.New("invalid book parameters")
		c.JSON(http.StatusBadRequest, &gin.Error{
			Err:  err,
			Type: http.StatusBadRequest,
			Meta: err.Error(),
		})
		return
	}

	err = b.bookUseCase.Create(bookJSON)
	if err != nil {
		c.JSON(http.StatusBadRequest, &gin.Error{
			Err:  err,
			Type: http.StatusBadRequest,
			Meta: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, &gin.H{
		"message": "book created",
	})
}
