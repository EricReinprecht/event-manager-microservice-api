package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/reinp/event-platform/backend/internal/appErrors"
	"github.com/reinp/event-platform/backend/internal/requests"
	"github.com/reinp/event-platform/backend/internal/service"
)

type PurchaseHandler struct {
	service *service.PurchaseService
}

func NewPurchaseHandler(
	service *service.PurchaseService,
) *PurchaseHandler {

	return &PurchaseHandler{
		service: service,
	}
}

func (h *PurchaseHandler) Create(c *gin.Context) {

	partyID, err := uuid.Parse(
		c.Param("id"),
	)

	if err != nil {

		c.JSON(
			http.StatusBadRequest,
			gin.H{"error": "invalid party id"},
		)

		return
	}

	userID, err := uuid.Parse(
		c.MustGet("user_id").(string),
	)

	if err != nil {

		c.JSON(
			http.StatusUnauthorized,
			gin.H{"error": "invalid user"},
		)

		return
	}

	var req requests.PurchaseRequest

	if err := c.ShouldBindJSON(&req); err != nil {

		c.JSON(
			http.StatusBadRequest,
			gin.H{"error": err.Error()},
		)

		return
	}

	purchase, err := h.service.CreatePurchase(
		c.Request.Context(),
		userID,
		partyID,
		req.Items,
	)

	if err != nil {

		switch {

		case errors.Is(err, appErrors.ErrTicketCategoryNotFound):

			c.JSON(
				http.StatusNotFound,
				gin.H{"error": err.Error()},
			)

		default:

			c.JSON(
				http.StatusBadRequest,
				gin.H{"error": err.Error()},
			)
		}

		return
	}

	c.JSON(
		http.StatusCreated,
		purchase,
	)
}
