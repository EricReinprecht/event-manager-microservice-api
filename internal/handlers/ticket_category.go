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

type updateTicketCategoryRequest struct {
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

func (h *TicketCategoryHandler) GetAll(c *gin.Context) {

	partyID, err := uuid.Parse(c.Param("id"))

	if err != nil {
		c.JSON(400, gin.H{
			"error": "invalid party id",
		})
		return
	}

	categories, err := h.service.FindByParty(
		c.Request.Context(),
		partyID,
	)

	if err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(200, categories)
}

func (h *TicketCategoryHandler) GetByID(c *gin.Context) {

	id, err := uuid.Parse(c.Param("id"))

	if err != nil {
		c.JSON(400, gin.H{
			"error": "invalid id",
		})
		return
	}

	category, err := h.service.FindByID(
		c.Request.Context(),
		id,
	)

	if err != nil {
		c.JSON(404, gin.H{
			"error": "ticket category not found",
		})
		return
	}

	c.JSON(200, category)
}

func (h *TicketCategoryHandler) Update(c *gin.Context) {

	id, err := uuid.Parse(c.Param("id"))

	if err != nil {
		c.JSON(400, gin.H{
			"error": "invalid id",
		})
		return
	}

	category, err := h.service.FindByID(
		c.Request.Context(),
		id,
	)

	if err != nil {
		c.JSON(404, gin.H{
			"error": "ticket category not found",
		})
		return
	}

	userID := uuid.MustParse(
		c.MustGet("user_id").(string),
	)

	_, err = h.partyService.FindOwnedParty(
		c.Request.Context(),
		category.PartyID,
		userID,
	)

	if err != nil {
		c.JSON(403, gin.H{
			"error": "not allowed",
		})
		return
	}

	var req updateTicketCategoryRequest

	if err := c.ShouldBindJSON(&req); err != nil {

		c.JSON(400, gin.H{
			"error": err.Error(),
		})

		return
	}

	category.Name = req.Name
	category.Price = req.Price
	category.Amount = req.Amount

	err = h.service.Update(
		c.Request.Context(),
		category,
	)

	if err != nil {

		if errors.Is(err, service.ErrTicketCategoryExists) {

			c.JSON(409, gin.H{
				"error": err.Error(),
			})

			return
		}

		c.JSON(500, gin.H{
			"error": err.Error(),
		})

		return
	}

	c.JSON(200, category)
}

func (h *TicketCategoryHandler) Delete(c *gin.Context) {

	id, err := uuid.Parse(c.Param("id"))

	if err != nil {
		c.JSON(400, gin.H{
			"error": "invalid id",
		})
		return
	}

	category, err := h.service.FindByID(
		c.Request.Context(),
		id,
	)

	if err != nil {
		c.JSON(404, gin.H{
			"error": "ticket category not found",
		})
		return
	}

	userID := uuid.MustParse(
		c.MustGet("user_id").(string),
	)

	_, err = h.partyService.FindOwnedParty(
		c.Request.Context(),
		category.PartyID,
		userID,
	)

	if err != nil {
		c.JSON(403, gin.H{
			"error": "not allowed",
		})
		return
	}

	err = h.service.Delete(
		c.Request.Context(),
		category,
	)

	if err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"message": "ticket category deleted",
	})
}
