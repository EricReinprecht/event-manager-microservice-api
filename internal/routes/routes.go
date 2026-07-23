package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/reinp/event-platform/backend/internal/handlers"
	"github.com/reinp/event-platform/backend/internal/middleware"
	"github.com/reinp/event-platform/backend/internal/service"
)

func Register(
	router *gin.Engine,
	authService *service.AuthService,
) {

	authHandler := handlers.NewAuthHandler(authService)

	// Public routes
	router.GET(
		"/health",
		handlers.Health,
	)

	router.POST(
		"/api/auth/register",
		authHandler.Register,
	)

	router.POST(
		"/api/auth/login",
		authHandler.Login,
	)

	// Protected routes
	protected := router.Group("/api")

	protected.Use(
		middleware.Auth(
			authService.Secret(),
		),
	)

	protected.GET(
		"/users/me",
		handlers.Me,
	)
}
