package helpers

import (
	"github.com/reinp/event-platform/backend/internal/payment/paypal"
	"github.com/reinp/event-platform/backend/internal/service"
)

func NewPaymentService(
	purchaseService *service.PurchaseService,
	ticketService *service.TicketService,
	paypalClient *paypal.Client,
) *service.PaymentService {

	return service.NewPaymentService(
		purchaseService,
		ticketService,
		paypalClient,
	)
}
