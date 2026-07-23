package server

import (
	"github.com/gin-gonic/gin"

	"github.com/reinp/event-platform/backend/internal/routes"
	"github.com/reinp/event-platform/backend/internal/service"
)

func Start(
	port string,
	authService *service.AuthService,
	partyService *service.PartyService,
) error {

	router := gin.Default()

	routes.Register(
		router,
		authService,
		partyService,
	)

	return router.Run(port)
}
