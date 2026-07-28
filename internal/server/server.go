package server

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"github.com/reinp/event-platform/backend/internal/middleware"
	"github.com/reinp/event-platform/backend/internal/routes"
	"github.com/reinp/event-platform/backend/internal/service"
)

func Start(
	port string,
	corsOrigins []string,
	authService *service.AuthService,
	userService *service.UserService,
	partyService *service.PartyService,
	categoryService *service.CategoryService,
	mediaService *service.MediaService,
	ticketCategoryService *service.TicketCategoryService,
	ticketService *service.TicketService,
	purchaseService *service.PurchaseService,
	paymentService *service.PaymentService,
	partyMemberService *service.PartyMemberService,
) error {

	globalLimiter := middleware.NewIPRateLimiter(
		20,
		50,
	)

	authLimiter := middleware.NewIPRateLimiter(
		1.0/12,
		5,
	)

	router := gin.Default()

	router.Use(
		cors.New(
			middleware.CORS(corsOrigins),
		),
	)

	router.Use(
		globalLimiter.Middleware(),
	)

	routes.Register(
		router,
		authLimiter,
		authService,
		userService,
		partyService,
		categoryService,
		mediaService,
		ticketCategoryService,
		ticketService,
		purchaseService,
		paymentService,
		partyMemberService,
	)

	return router.Run(port)
}
