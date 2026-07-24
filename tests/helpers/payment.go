package helpers

import (
	"github.com/reinp/event-platform/backend/internal/database"
	"github.com/reinp/event-platform/backend/internal/payment/paypal"
	"github.com/reinp/event-platform/backend/internal/repository"
	"github.com/reinp/event-platform/backend/internal/service"
)

func NewPaymentService(
	executor database.DBExecutor,
	purchaseService *service.PurchaseService,
	ticketService *service.TicketService,
	paypalClient *paypal.Client,
) *service.PaymentService {

	paymentEventRepository := repository.NewPaymentEventRepository(
		executor,
	)

	return service.NewPaymentService(
		purchaseService,
		ticketService,
		paypalClient,
		paymentEventRepository,
	)
}
