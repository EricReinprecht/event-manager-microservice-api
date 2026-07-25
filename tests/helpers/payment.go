package helpers

import (
	"github.com/reinp/event-platform/backend/internal/database"
	"github.com/reinp/event-platform/backend/internal/payment"
	"github.com/reinp/event-platform/backend/internal/repository"
	"github.com/reinp/event-platform/backend/internal/service"
)

func NewPaymentService(
	executor database.DBExecutor,
	purchaseService *service.PurchaseService,
	ticketService *service.TicketService,
	paymentGateway payment.Gateway,
) *service.PaymentService {

	paymentEventRepository := repository.NewPaymentEventRepository(
		executor,
	)

	purchaseRepository := repository.NewPurchaseRepository(
		executor,
	)

	partyMemberRepository := repository.NewPartyMemberRepository(
		executor,
	)

	partyRepository := repository.NewPartyRepository(
		executor,
	)

	roleRepository := repository.NewPartyMemberRoleRepository(
		executor,
	)

	partyMemberService := service.NewPartyMemberService(
		partyMemberRepository,
		partyRepository,
		roleRepository,
	)

	return service.NewPaymentService(
		purchaseService,
		ticketService,
		partyMemberService,
		paymentGateway,
		paymentEventRepository,
		purchaseRepository,
	)
}
