package responses

import (
	"math"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Pagination struct {
	Data       any   `json:"data"`
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"totalPages"`
}

func Error(
	c *gin.Context,
	status int,
	err error,
) {
	c.JSON(
		status,
		gin.H{
			"error": err.Error(),
		},
	)
}

func Success(
	c *gin.Context,
	status int,
	data any,
) {
	c.JSON(
		status,
		data,
	)
}

func Paginated(
	c *gin.Context,
	data any,
	page int,
	limit int,
	total int64,
) {

	totalPages := 0

	if limit > 0 {
		totalPages = int(
			math.Ceil(
				float64(total) /
					float64(limit),
			),
		)
	}

	c.JSON(
		http.StatusOK,
		Pagination{
			Data:       data,
			Page:       page,
			Limit:      limit,
			Total:      total,
			TotalPages: totalPages,
		},
	)
}
