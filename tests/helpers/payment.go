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

	return service.NewPaymentService(
		purchaseService,
		ticketService,
		paymentGateway,
		paymentEventRepository,
		purchaseRepository,
	)
}
