package helpers

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

func QueryPagination(c *gin.Context) (int, int) {

	page, _ := strconv.Atoi(
		c.DefaultQuery("page", "1"),
	)

	limit, _ := strconv.Atoi(
		c.DefaultQuery("limit", "10"),
	)

	if page < 1 {
		page = 1
	}

	if limit < 1 || limit > 100 {
		limit = 10
	}

	return page, limit
}
