package helpers

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func UUIDParam(
	c *gin.Context,
	key string,
) (uuid.UUID, error) {

	id, err := uuid.Parse(
		c.Param(key),
	)

	if err != nil {
		return uuid.Nil, errors.New("invalid id")
	}

	return id, nil
}
