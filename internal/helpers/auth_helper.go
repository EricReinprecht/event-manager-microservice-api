package helpers

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/reinp/event-platform/backend/internal/responses"
)

func MustUserID(c *gin.Context) (uuid.UUID, bool) {
	userID, ok := RequireUserID(c)

	if !ok {
		responses.Unauthorized(c)
		return uuid.Nil, false
	}

	return userID, true
}
