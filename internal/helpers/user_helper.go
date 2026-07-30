package helpers

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func RequireUserID(
	c *gin.Context,
) (uuid.UUID, bool) {

	value, exists := c.Get("userID")

	if !exists {
		return uuid.Nil, false
	}

	switch id := value.(type) {

	case uuid.UUID:
		return id, true

	case string:
		parsed, err := uuid.Parse(id)

		if err != nil {
			return uuid.Nil, false
		}

		return parsed, true
	}

	return uuid.Nil, false
}
