package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/reinp/event-platform/backend/internal/models"
	"github.com/reinp/event-platform/backend/internal/service"
)

type PartyHandler struct {
	service *service.PartyService
}

func NewPartyHandler(
	service *service.PartyService,
) *PartyHandler {

	return &PartyHandler{
		service: service,
	}
}

type createPartyRequest struct {
	Title string `json:"title"`

	Description string `json:"description"`

	Location string `json:"location"`

	StartAt time.Time `json:"start_at" binding:"required"`

	EndAt time.Time `json:"end_at" binding:"required"`

	CategoryID uuid.UUID `json:"category_id"`

	ThumbnailID *uuid.UUID `json:"thumbnail_id"`
}

type updatePartyRequest struct {
	Title string `json:"title"`

	Description string `json:"description"`

	Location string `json:"location"`

	StartAt time.Time `json:"start_at" binding:"required"`

	EndAt time.Time `json:"end_at" binding:"required"`

	CategoryID uuid.UUID `json:"category_id"`

	ThumbnailID *uuid.UUID `json:"thumbnail_id"`
}

func (h *PartyHandler) Create(c *gin.Context) {

	var req createPartyRequest

	if err := c.ShouldBindJSON(&req); err != nil {

		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})

		return
	}

	if req.EndAt.Before(req.StartAt) {

		c.JSON(http.StatusBadRequest, gin.H{
			"error": "end date must be after start date",
		})

		return
	}

	userID, exists := c.Get("user_id")

	if !exists {

		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "not authenticated",
		})

		return
	}

	organizerID, err := uuid.Parse(
		userID.(string),
	)

	if err != nil {

		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid user id",
		})

		return
	}

	party := &models.Party{

		Title: req.Title,

		Description: req.Description,

		Location: req.Location,

		CategoryID: req.CategoryID,

		OrganizerID: organizerID,

		ThumbnailID: req.ThumbnailID,
	}

	err = h.service.Create(
		c.Request.Context(),
		party,
	)

	if err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})

		return
	}

	c.JSON(http.StatusCreated, party)
}

func (h *PartyHandler) GetAll(c *gin.Context) {

	parties, err := h.service.FindAll(
		c.Request.Context(),
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(
		http.StatusOK,
		parties,
	)
}

func (h *PartyHandler) GetByID(c *gin.Context) {

	id, err := uuid.Parse(
		c.Param("id"),
	)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid id",
		})
		return
	}

	party, err := h.service.FindByID(
		c.Request.Context(),
		id,
	)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "party not found",
		})
		return
	}

	c.JSON(
		http.StatusOK,
		party,
	)
}

func (h *PartyHandler) Update(c *gin.Context) {

	id, err := uuid.Parse(
		c.Param("id"),
	)

	if err != nil {
		c.JSON(400, gin.H{
			"error": "invalid id",
		})
		return
	}

	party, err := h.service.FindByID(
		c.Request.Context(),
		id,
	)

	if err != nil {
		c.JSON(404, gin.H{
			"error": "party not found",
		})
		return
	}

	userID := c.MustGet("user_id").(string)

	if party.OrganizerID.String() != userID {

		c.JSON(403, gin.H{
			"error": "not allowed",
		})

		return
	}

	var req updatePartyRequest

	if err := c.ShouldBindJSON(&req); err != nil {

		if req.EndAt.Before(req.StartAt) {

			c.JSON(http.StatusBadRequest, gin.H{
				"error": "end date must be after start date",
			})

			return
		}

		c.JSON(400, gin.H{
			"error": err.Error(),
		})

		return
	}

	party.Title = req.Title

	party.Description = req.Description

	party.Location = req.Location

	party.CategoryID = req.CategoryID

	party.ThumbnailID = req.ThumbnailID

	err = h.service.Update(
		c.Request.Context(),
		party,
	)

	if err != nil {

		c.JSON(500, gin.H{
			"error": err.Error(),
		})

		return
	}

	c.JSON(200, party)
}

func (h *PartyHandler) Delete(c *gin.Context) {

	id, err := uuid.Parse(
		c.Param("id"),
	)

	if err != nil {
		c.JSON(400, gin.H{
			"error": "invalid id",
		})
		return
	}

	party, err := h.service.FindByID(
		c.Request.Context(),
		id,
	)

	if err != nil {
		c.JSON(404, gin.H{
			"error": "party not found",
		})
		return
	}

	userID := c.MustGet("user_id").(string)

	if party.OrganizerID.String() != userID {

		c.JSON(403, gin.H{
			"error": "not allowed",
		})

		return
	}

	err = h.service.Delete(
		c.Request.Context(),
		party,
	)

	if err != nil {

		c.JSON(500, gin.H{
			"error": err.Error(),
		})

		return
	}

	c.JSON(200, gin.H{
		"message": "party deleted",
	})
}
