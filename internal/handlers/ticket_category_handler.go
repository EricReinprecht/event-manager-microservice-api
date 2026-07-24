package handlers

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/reinp/event-platform/backend/internal/appErrors"
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

	UnitPrice int64 `json:"price" binding:"required"`

	Capacity int `json:"capacity" binding:"required"`

	RequiresVerification bool `json:"requires_verification"`

	AccessWindows []createTicketAccessWindowRequest `json:"access_windows" binding:"required"`
}

type updateTicketCategoryRequest struct {
	Name string `json:"name" binding:"required"`

	UnitPrice int64 `json:"price" binding:"required"`

	Capacity int `json:"capacity" binding:"required"`

	RequiresVerification bool `json:"requires_verification"`

	AccessWindows []createTicketAccessWindowRequest `json:"access_windows" binding:"required"`
}

type createTicketAccessWindowRequest struct {
	StartsAt time.Time `json:"starts_at" binding:"required"`

	EndsAt time.Time `json:"ends_at" binding:"required"`
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

	if len(req.AccessWindows) == 0 {

		c.JSON(http.StatusBadRequest, gin.H{
			"error": "ticket category requires at least one access window",
		})

		return
	}

	category := &models.TicketCategory{

		Name: req.Name,

		Price: req.UnitPrice,

		Capacity: req.Capacity,

		RequiresVerification: req.RequiresVerification,

		PartyID: partyID,
	}

	for _, window := range req.AccessWindows {

		if window.EndsAt.Before(window.StartsAt) {

			c.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid access window: end time must be later than start time",
			})

			return
		}

		category.AccessWindows = append(
			category.AccessWindows,
			models.TicketAccessWindow{
				StartsAt: window.StartsAt,
				EndsAt:   window.EndsAt,
			},
		)
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

		if errors.Is(err, appErrors.ErrTicketAccessWindowRequired) {

			c.JSON(http.StatusBadRequest, gin.H{
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
	category.Price = req.UnitPrice
	category.Capacity = req.Capacity
	category.RequiresVerification = req.RequiresVerification
	category.AccessWindows = nil

	for _, window := range req.AccessWindows {

		if window.EndsAt.Before(window.StartsAt) {

			c.JSON(http.StatusBadRequest, gin.H{
				"error": "access window end must be after start",
			})

			return
		}

		category.AccessWindows = append(
			category.AccessWindows,
			models.TicketAccessWindow{
				StartsAt:         window.StartsAt,
				EndsAt:           window.EndsAt,
				TicketCategoryID: category.ID,
			},
		)
	}

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
