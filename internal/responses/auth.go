package responses

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/reinp/event-platform/backend/internal/appErrors"
)

func Forbidden(
	c *gin.Context,
) {
	Error(
		c,
		http.StatusForbidden,
		appErrors.ErrNotAllowed,
	)
}
