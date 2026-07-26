package helpers

import (
	"github.com/reinp/event-platform/backend/internal/database"
	"github.com/reinp/event-platform/backend/internal/payment"
	"github.com/reinp/event-platform/backend/internal/repository"
	"github.com/reinp/event-platform/backend/internal/service"
	"gorm.io/gorm"
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

	refundService := service.NewRefundService()

	return service.NewPaymentService(
		purchaseService,
		ticketService,
		partyMemberService,
		paymentGateway,
		paymentEventRepository,
		purchaseRepository,
		refundService,
	)
}

func SetupPaymentTestService(
	db *gorm.DB,
	gateway payment.Gateway,
) *service.PaymentService {

	executor := database.NewGormExecutor(db)

	return NewPaymentService(
		executor,
		NewPurchaseService(db),
		NewTicketService(db),
		gateway,
	)
}
