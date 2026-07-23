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
	categoryService *service.CategoryService,
	mediaService *service.MediaService,
) error {

	router := gin.Default()

	routes.Register(
		router,
		authService,
		partyService,
		categoryService,
		mediaService,
	)

	return router.Run(port)
}
