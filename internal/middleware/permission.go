package middleware

import (
	"github.com/gin-gonic/gin"

	"github.com/reinp/event-platform/backend/internal/helpers"
	"github.com/reinp/event-platform/backend/internal/models/enum"
	"github.com/reinp/event-platform/backend/internal/responses"
	"github.com/reinp/event-platform/backend/internal/service"
)

func PartyOwnerMiddleware(
	permission *service.PermissionService,
) gin.HandlerFunc {

	return func(c *gin.Context) {

		partyID, err := helpers.UUIDParam(
			c,
			"id",
		)

		if err != nil {

			responses.BadRequest(
				c,
				err,
			)

			c.Abort()
			return
		}

		userID, ok := helpers.RequireUserID(c)

		if !ok {

			responses.Unauthorized(c)

			c.Abort()
			return
		}

		err = permission.RequirePartyRole(
			c.Request.Context(),
			partyID,
			userID,
			enum.RoleOrganizer,
			enum.RoleAdmin,
		)

		if err != nil {

			responses.HandleDomainError(
				c,
				err,
			)

			c.Abort()
			return
		}

		c.Set(
			"partyID",
			partyID,
		)

		c.Next()
	}
}
