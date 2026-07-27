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

func (h *UserHandler) Me(
	c *gin.Context,
) {

	userIDValue, exists := c.Get("user_id")

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
