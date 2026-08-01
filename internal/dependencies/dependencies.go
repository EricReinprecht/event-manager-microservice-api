package dependencies

import (
	"github.com/reinp/event-platform/backend/internal/service"
)

type Container struct {
	AuthService *service.AuthService

	UserService *service.UserService

	PartyService *service.PartyService

	CategoryService *service.CategoryService

	MediaService *service.MediaService

	TicketCategoryService *service.TicketCategoryService

	TicketService *service.TicketService

	PurchaseService *service.PurchaseService

	PaymentService *service.PaymentService

	PartyMemberService *service.PartyMemberService

	PermissionService *service.PermissionService
}
