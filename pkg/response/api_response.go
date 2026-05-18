package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type ApiResponse[T any] struct {
	StatusCode int    `json:"statusCode"`
	HttpStatus string `json:"httpStatus"`
	Message    string `json:"message"`
	Data       T      `json:"data"`
}

func OK(c *gin.Context, data any) {
	c.JSON(http.StatusOK, ApiResponse[any]{
		StatusCode: http.StatusOK,
		HttpStatus: "OK",
		Message:    "OK",
		Data:       data,
	})
}

func Created(c *gin.Context, data any) {
	c.JSON(http.StatusCreated, ApiResponse[any]{
		StatusCode: http.StatusCreated,
		HttpStatus: "CREATED",
		Message:    "CREATED",
		Data:       data,
	})
}

func BadRequest(c *gin.Context, message string) {
	c.JSON(http.StatusBadRequest, ApiResponse[any]{
		StatusCode: http.StatusBadRequest,
		HttpStatus: "BAD_REQUEST",
		Message:    message,
		Data:       nil,
	})
}

func Unauthorized(c *gin.Context, message string) {
	c.JSON(http.StatusUnauthorized, ApiResponse[any]{
		StatusCode: http.StatusUnauthorized,
		HttpStatus: "UNAUTHORIZED",
		Message:    message,
		Data:       nil,
	})
}

func InternalServerError(c *gin.Context, message string) {
	c.JSON(http.StatusInternalServerError, ApiResponse[any]{
		StatusCode: http.StatusInternalServerError,
		HttpStatus: "INTERNAL_SERVER_ERROR",
		Message:    message,
		Data:       nil,
	})
}
