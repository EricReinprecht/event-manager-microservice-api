package handlers

import (
	"net/http"

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
	Title string `json:"title" binding:"required"`

	Description string `json:"description"`

	Date string `json:"date" binding:"required"`

	Location string `json:"location"`
}

func (h *PartyHandler) Create(c *gin.Context) {

	var req createPartyRequest

	if err := c.ShouldBindJSON(&req); err != nil {

		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
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

		OrganizerID: organizerID,
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
