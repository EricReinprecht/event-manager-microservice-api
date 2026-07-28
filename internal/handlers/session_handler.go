package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/reinp/event-platform/backend/internal/service"
)

type SessionHandler struct {
	service *service.AuthService
}

func NewSessionHandler(
	service *service.AuthService,
) *SessionHandler {

	return &SessionHandler{
		service: service,
	}
}

func (h *SessionHandler) GetSessions(c *gin.Context) {

	userID := c.MustGet(
		"userID",
	).(uuid.UUID)

	familyID := c.MustGet(
		"familyID",
	).(uuid.UUID)

	sessions, err := h.service.Sessions(
		c.Request.Context(),
		userID,
		familyID,
	)

	if err != nil {
		c.JSON(
			500,
			gin.H{
				"error": err.Error(),
			},
		)

		return
	}

	c.JSON(
		http.StatusOK,
		sessions,
	)
}

func (h *SessionHandler) DeleteSession(
	c *gin.Context,
) {

	userID := c.MustGet(
		"userID",
	).(uuid.UUID)

	familyID, err := uuid.Parse(
		c.Param("familyID"),
	)

	if err != nil {

		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": "invalid session id",
			},
		)

		return
	}

	err = h.service.RevokeSession(
		c.Request.Context(),
		userID,
		familyID,
	)

	if err != nil {

		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": err.Error(),
			},
		)

		return
	}

	c.JSON(
		http.StatusOK,
		gin.H{
			"message": "session revoked",
		},
	)
}

func (h *SessionHandler) LogoutAll(
	c *gin.Context,
) {

	userID := c.MustGet(
		"userID",
	).(uuid.UUID)

	err := h.service.LogoutAll(
		c.Request.Context(),
		userID,
	)

	if err != nil {

		c.JSON(
			http.StatusInternalServerError,
			gin.H{
				"error": err.Error(),
			},
		)

		return
	}

	c.JSON(
		http.StatusOK,
		gin.H{
			"message": "logged out from all devices",
		},
	)
}
