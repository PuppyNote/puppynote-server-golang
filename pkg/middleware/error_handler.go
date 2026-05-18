package middleware

import (
	"net/http"

	pnerrors "github.com/PuppyNote/puppynote-server-golang/pkg/errors"
	"github.com/PuppyNote/puppynote-server-golang/pkg/response"
	"github.com/gin-gonic/gin"
)

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) == 0 {
			return
		}

		err := c.Errors.Last().Err
		switch e := err.(type) {
		case *pnerrors.PuppyNoteError:
			switch e.Code {
			case http.StatusUnauthorized:
				response.Unauthorized(c, e.Message)
			default:
				response.BadRequest(c, e.Message)
			}
		default:
			response.InternalServerError(c, "서버 오류가 발생했습니다.")
		}
	}
}
