package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	appErrors "github.com/reinp/event-platform/backend/internal/appErrors"
	"github.com/reinp/event-platform/backend/internal/models"
	"github.com/reinp/event-platform/backend/internal/models/enum"
	"github.com/reinp/event-platform/backend/internal/service"
)

type PartyMemberHandler struct {
	service *service.PartyMemberService
}

func NewPartyMemberHandler(
	service *service.PartyMemberService,
) *PartyMemberHandler {

	return &PartyMemberHandler{
		service: service,
	}
}

type createPartyMemberRequest struct {
	UserID uuid.UUID      `json:"user_id" binding:"required"`
	Role   enum.PartyRole `json:"role" binding:"required"`
}

type updatePartyMemberRolesRequest struct {
	Roles []enum.PartyRole `json:"roles" binding:"required"`
}

func (h *PartyMemberHandler) Create(
	c *gin.Context,
) {

	partyID, err := uuid.Parse(
		c.Param("id"),
	)

	if err != nil {
		c.JSON(400, gin.H{
			"error": "invalid party id",
		})
		return
	}

	currentUser, err := uuid.Parse(
		c.MustGet("user_id").(string),
	)

	if err != nil {
		c.JSON(401, gin.H{
			"error": "invalid user",
		})
		return
	}

	if !h.service.HasRole(
		c.Request.Context(),
		partyID,
		currentUser,
		enum.RoleOrganizer,
		enum.RoleAdmin,
	) {

		c.JSON(403, gin.H{
			"error": "not allowed",
		})

		return
	}

	var req createPartyMemberRequest

	if err := c.ShouldBindJSON(&req); err != nil {

		c.JSON(400, gin.H{
			"error": err.Error(),
		})

		return
	}

	member := models.PartyMember{
		UserID:  req.UserID,
		PartyID: partyID,

		Roles: []models.PartyMemberRole{
			{
				ID:   uuid.New(),
				Role: enum.RoleOrganizer,
			},
		},
	}

	if err != nil {

		switch {

		case errors.Is(err, appErrors.ErrPartyMemberAlreadyExists):

			c.JSON(409, gin.H{
				"error": err.Error(),
			})

		case errors.Is(err, appErrors.ErrInvalidPartyMemberRole):

			c.JSON(400, gin.H{
				"error": err.Error(),
			})

		default:

			c.JSON(500, gin.H{
				"error": err.Error(),
			})
		}

		return
	}
	c.JSON(http.StatusCreated, member)
}

func (h *PartyMemberHandler) GetAll(
	c *gin.Context,
) {

	partyID, err := uuid.Parse(
		c.Param("id"),
	)

	if err != nil {

		c.JSON(400, gin.H{
			"error": "invalid party id",
		})

		return
	}

	members, err := h.service.FindByParty(
		c.Request.Context(),
		partyID,
	)

	if err != nil {

		c.JSON(500, gin.H{
			"error": err.Error(),
		})

		return
	}

	c.JSON(
		200,
		members,
	)
}

func (h *PartyMemberHandler) Delete(
	c *gin.Context,
) {

	memberID, err := uuid.Parse(
		c.Param("memberID"),
	)

	if err != nil {

		c.JSON(400, gin.H{
			"error": "invalid member id",
		})

		return
	}

	err = h.service.Delete(
		c.Request.Context(),
		memberID,
	)

	if err != nil {

		switch {

		case errors.Is(err, appErrors.ErrCannotRemoveOrganizer):

			c.JSON(403, gin.H{
				"error": err.Error(),
			})

		case errors.Is(err, appErrors.ErrPartyMemberNotFound):

			c.JSON(404, gin.H{
				"error": err.Error(),
			})

		default:

			c.JSON(500, gin.H{
				"error": err.Error(),
			})
		}

		return
	}

	c.JSON(200, gin.H{
		"message": "member removed",
	})
}

func (h *PartyMemberHandler) UpdateRoles(
	c *gin.Context,
) {

	memberID, err := uuid.Parse(
		c.Param("memberID"),
	)

	if err != nil {

		c.JSON(400, gin.H{
			"error": "invalid member id",
		})

		return
	}

	var req updatePartyMemberRolesRequest

	if err := c.ShouldBindJSON(&req); err != nil {

		c.JSON(400, gin.H{
			"error": err.Error(),
		})

		return
	}

	err = h.service.SyncRoles(
		c.Request.Context(),
		memberID,
		req.Roles,
	)

	if err != nil {

		c.JSON(500, gin.H{
			"error": err.Error(),
		})

		return
	}

	c.JSON(200, gin.H{
		"message": "roles updated",
	})
}
