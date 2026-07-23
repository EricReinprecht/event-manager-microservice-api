package server

import (
	"github.com/gin-gonic/gin"

	"github.com/reinp/event-platform/backend/internal/routes"
	"github.com/reinp/event-platform/backend/internal/service"
)

func Start(
	addr string,
	authService *service.AuthService,
) error {

	router := gin.Default()

	routes.Register(
		router,
		authService,
	)

	return router.Run(addr)
}
