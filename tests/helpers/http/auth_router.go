package http

import (
	"github.com/gin-gonic/gin"

	"github.com/reinp/event-platform/backend/internal/handlers"
)

func AuthRouter(
	handler *handlers.SessionHandler,
	userID interface{},
	familyID interface{},
) *gin.Engine {

	router := gin.New()

	router.GET(
		"/api/auth/sessions",
		func(c *gin.Context) {

			c.Set(
				"userID",
				userID,
			)

			c.Set(
				"familyID",
				familyID,
			)

			handler.GetSessions(c)
		},
	)

	router.DELETE(
		"/api/auth/sessions/:familyID",
		func(c *gin.Context) {

			c.Set(
				"userID",
				userID,
			)

			handler.DeleteSession(c)
		},
	)

	router.DELETE(
		"/api/auth/sessions",
		func(c *gin.Context) {

			c.Set(
				"userID",
				userID,
			)

			handler.LogoutAll(c)
		},
	)

	return router
}
