package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/reinp/event-platform/backend/internal/helpers"
	"github.com/reinp/event-platform/backend/internal/responses"
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

func (h *SessionHandler) GetSessions(
	c *gin.Context,
) {

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

		responses.HandleDomainError(
			c,
			err,
		)

		return
	}

	responses.Success(
		c,
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

	familyID, err := helpers.UUIDParam(
		c,
		"familyID",
	)

	if err != nil {

		responses.BadRequest(
			c,
			err,
		)

		return
	}

	err = h.service.RevokeSession(
		c.Request.Context(),
		userID,
		familyID,
	)

	if err != nil {

		responses.HandleDomainError(
			c,
			err,
		)

		return
	}

	responses.Success(
		c,
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

		responses.HandleDomainError(
			c,
			err,
		)

		return
	}

	responses.Success(
		c,
		http.StatusOK,
		gin.H{
			"message": "logged out from all devices",
		},
	)
}
