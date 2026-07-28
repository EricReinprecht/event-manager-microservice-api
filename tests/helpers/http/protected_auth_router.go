package http

import (
	"github.com/gin-gonic/gin"

	"github.com/reinp/event-platform/backend/internal/middleware"
	"github.com/reinp/event-platform/backend/internal/service"
)

func ProtectedAuthRouter(
	authService *service.AuthService,
) *gin.Engine {

	router := gin.New()

	router.Use(
		middleware.Auth(
			authService,
		),
	)

	router.GET(
		"/api/auth/sessions",
		func(c *gin.Context) {
			c.JSON(
				200,
				gin.H{
					"message": "authorized",
				},
			)
		},
	)

	router.DELETE(
		"/api/auth/sessions",
		func(c *gin.Context) {
			c.JSON(
				200,
				gin.H{
					"message": "authorized",
				},
			)
		},
	)

	return router
}
