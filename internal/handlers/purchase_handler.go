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
			gin.H{
				"error": "invalid party id",
			},
		)

		return
	}

	userID, err := uuid.Parse(
		c.MustGet("user_id").(string),
	)

	if err != nil {

		c.JSON(
			http.StatusUnauthorized,
			gin.H{
				"error": "invalid user",
			},
		)

		return
	}

	var req requests.PurchaseRequest

	if err := c.ShouldBindJSON(&req); err != nil {

		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": err.Error(),
			},
		)

		return
	}

	if len(req.Items) == 0 {

		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": "purchase requires items",
			},
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

		case errors.Is(
			err,
			appErrors.ErrPartyNotFound,
		):

			c.JSON(
				http.StatusNotFound,
				gin.H{
					"error": err.Error(),
				},
			)

		case errors.Is(
			err,
			appErrors.ErrTicketCategoryNotFound,
		):

			c.JSON(
				http.StatusNotFound,
				gin.H{
					"error": err.Error(),
				},
			)

		case errors.Is(
			err,
			appErrors.ErrNotEnoughTickets,
		):

			c.JSON(
				http.StatusConflict,
				gin.H{
					"error": err.Error(),
				},
			)

		default:

			c.JSON(
				http.StatusBadRequest,
				gin.H{
					"error": err.Error(),
				},
			)
		}

		return
	}

	c.JSON(
		http.StatusCreated,
		purchase,
	)
}

func (h *PurchaseHandler) GetByID(c *gin.Context) {

	id, err := uuid.Parse(
		c.Param("id"),
	)

	if err != nil {

		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": "invalid purchase id",
			},
		)

		return
	}

	purchase, err := h.service.GetPurchase(
		c.Request.Context(),
		id,
	)

	if err != nil {

		c.JSON(
			http.StatusNotFound,
			gin.H{
				"error": err.Error(),
			},
		)

		return
	}

	c.JSON(
		http.StatusOK,
		purchase,
	)
}
