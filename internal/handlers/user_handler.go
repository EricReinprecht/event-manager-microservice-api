package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/reinp/event-platform/backend/internal/service"
)

type UserHandler struct {
	service *service.UserService
}

func NewUserHandler(
	service *service.UserService,
) *UserHandler {

	return &UserHandler{
		service: service,
	}
}

type CompleteProfileRequest struct {
	FirstName string `json:"firstName" binding:"required"`
	LastName  string `json:"lastName" binding:"required"`
}

type ChangePasswordRequest struct {
	OldPassword string `json:"oldPassword" binding:"required"`
	NewPassword string `json:"newPassword" binding:"required"`
}

func (h *UserHandler) Me(
	c *gin.Context,
) {

	userIDValue, exists := c.Get("userID")

	if !exists {

		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "not authenticated",
		})

		return
	}

	userID, ok := userIDValue.(uuid.UUID)

	if !ok {

		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "invalid user id",
		})

		return
	}

	user, err := h.service.GetByID(
		c.Request.Context(),
		userID,
	)

	if err != nil {

		c.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{

		"id": user.ID,

		"email": user.Email,

		"username": user.Username,

		"first_name": user.FirstName,

		"last_name": user.LastName,

		"profile_completed": user.ProfileCompleted,
	})
}

func (h *UserHandler) CompleteProfile(
	c *gin.Context,
) {

	userIDValue, exists := c.Get("userID")

	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "not authenticated",
		})

		return
	}

	userID, ok := userIDValue.(uuid.UUID)

	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "invalid user id",
		})

		return
	}

	var req CompleteProfileRequest

	if err := c.ShouldBindJSON(&req); err != nil {

		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})

		return
	}

	user, err := h.service.CompleteProfile(
		c.Request.Context(),
		userID,
		req.FirstName,
		req.LastName,
	)

	if err != nil {

		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{

		"id": user.ID,

		"email": user.Email,

		"username": user.Username,

		"first_name": user.FirstName,

		"last_name": user.LastName,

		"profile_completed": true,
	})
}

func (h *UserHandler) ChangePassword(
	c *gin.Context,
) {

	userID := c.MustGet(
		"userID",
	).(uuid.UUID)

	var req ChangePasswordRequest

	if err := c.ShouldBindJSON(
		&req,
	); err != nil {

		c.JSON(
			400,
			gin.H{
				"error": err.Error(),
			},
		)

		return
	}

	err := h.service.ChangePassword(
		c.Request.Context(),
		userID,
		req.OldPassword,
		req.NewPassword,
	)

	if err != nil {

		c.JSON(
			400,
			gin.H{
				"error": err.Error(),
			},
		)

		return
	}

	c.JSON(
		200,
		gin.H{
			"message": "password changed",
		},
	)
}
