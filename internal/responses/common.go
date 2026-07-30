package responses

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

func InternalError(
	c *gin.Context,
	err error,
) {
	Error(
		c,
		http.StatusInternalServerError,
		err,
	)
}

func BadRequest(
	c *gin.Context,
	err error,
) {
	Error(
		c,
		http.StatusBadRequest,
		err,
	)
}

func Unauthorized(c *gin.Context) {
	Error(
		c,
		http.StatusUnauthorized,
		errors.New("unauthorized"),
	)
}

func NotFound(
	c *gin.Context,
	err error,
) {
	Error(
		c,
		http.StatusNotFound,
		err,
	)
}
