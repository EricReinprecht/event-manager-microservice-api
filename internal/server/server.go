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
	ticketCategoryService *service.TicketCategoryService,
	ticketService *service.TicketService,
	partyMemberService *service.PartyMemberService,
) error {

	router := gin.Default()

	routes.Register(
		router,
		authService,
		partyService,
		categoryService,
		mediaService,
		ticketCategoryService,
		ticketService,
		partyMemberService,
	)

	return router.Run(port)
}
