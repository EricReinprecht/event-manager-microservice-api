package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/reinp/event-platform/backend/internal/models"
	"github.com/reinp/event-platform/backend/internal/service"
)

type TicketCategoryHandler struct {
	service      *service.TicketCategoryService
	partyService *service.PartyService
}

func NewTicketCategoryHandler(
	service *service.TicketCategoryService,
	partyService *service.PartyService,
) *TicketCategoryHandler {

	return &TicketCategoryHandler{
		service:      service,
		partyService: partyService,
	}
}

type createTicketCategoryRequest struct {
	Name string `json:"name" binding:"required"`

	Price float64 `json:"price" binding:"required"`

	Amount int `json:"amount" binding:"required"`
}

func (h *TicketCategoryHandler) Create(c *gin.Context) {

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

	_, err = h.partyService.FindOwnedParty(
		c.Request.Context(),
		partyID,
		userID,
	)

	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "not allowed",
		})
		return
	}

	var req createTicketCategoryRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	category := &models.TicketCategory{
		Name:    req.Name,
		Price:   req.Price,
		Amount:  req.Amount,
		PartyID: partyID,
	}

	err = h.service.Create(
		c.Request.Context(),
		category,
	)

	if err != nil {

		if errors.Is(err, service.ErrTicketCategoryExists) {

			c.JSON(http.StatusConflict, gin.H{
				"error": err.Error(),
			})

			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})

		return
	}

	c.JSON(http.StatusCreated, category)
}
