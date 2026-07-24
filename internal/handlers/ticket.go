package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	appErrors "github.com/reinp/event-platform/backend/internal/appErrors"
	"github.com/reinp/event-platform/backend/internal/dto"
	"github.com/reinp/event-platform/backend/internal/service"
)

type TicketHandler struct {
	service *service.TicketService
}

func NewTicketHandler(
	service *service.TicketService,
) *TicketHandler {

	return &TicketHandler{
		service: service,
	}
}

type purchaseTicketRequest struct {
	Items []dto.PurchaseTicketItem `json:"items" binding:"required"`
}

type scanTicketRequest struct {
	Code string `json:"code" binding:"required"`
}

func (h *TicketHandler) Purchase(c *gin.Context) {

	partyID, err := uuid.Parse(
		c.Param("id"),
	)

	if err != nil {

		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid party id",
		})

		return
	}

	userID, err := uuid.Parse(
		c.MustGet("user_id").(string),
	)

	if err != nil {

		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "invalid user",
		})

		return
	}

	var req purchaseTicketRequest

	if err := c.ShouldBindJSON(&req); err != nil {

		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})

		return
	}

	if len(req.Items) == 0 {

		c.JSON(http.StatusBadRequest, gin.H{
			"error": "no tickets selected",
		})

		return
	}

	for _, item := range req.Items {

		if item.TicketCategoryID == uuid.Nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "ticket category id is required",
			})
			return
		}

		if item.Quantity <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "quantity must be greater than zero",
			})
			return
		}
	}

	purchase, err := h.service.CreatePurchase(
		c.Request.Context(),
		partyID,
		userID,
		req.Items,
	)

	if err != nil {

		switch {

		case errors.Is(err, appErrors.ErrPartyNotFound):

			c.JSON(http.StatusNotFound, gin.H{
				"error": err.Error(),
			})

		case errors.Is(err, appErrors.ErrTicketCategoryNotFound):

			c.JSON(http.StatusNotFound, gin.H{
				"error": err.Error(),
			})

		case errors.Is(err, appErrors.ErrNotEnoughTickets):

			c.JSON(http.StatusConflict, gin.H{
				"error": err.Error(),
			})

		default:

			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
		}

		return
	}

	c.JSON(http.StatusCreated, purchase)
}

func (h *TicketHandler) Scan(c *gin.Context) {

	var req scanTicketRequest

	if err := c.ShouldBindJSON(&req); err != nil {

		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})

		return
	}

	userID, err := uuid.Parse(
		c.MustGet("user_id").(string),
	)

	if err != nil {

		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "invalid user",
		})

		return
	}

	ticket, err := h.service.Scan(
		c.Request.Context(),
		userID,
		req.Code,
	)

	if err != nil {

		switch {

		case errors.Is(err, appErrors.ErrTicketAlreadyUsed):

			c.JSON(http.StatusConflict, gin.H{
				"error": err.Error(),
			})

		default:

			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
		}

		return
	}

	c.JSON(http.StatusOK, ticket)
}

func (h *TicketHandler) GetMyTickets(c *gin.Context) {

	userID, err := uuid.Parse(
		c.MustGet("user_id").(string),
	)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "invalid user",
		})

		return
	}

	tickets, err := h.service.GetMyTickets(
		c.Request.Context(),
		userID,
	)

	if err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})

		return
	}

	c.JSON(
		http.StatusOK,
		tickets,
	)
}
