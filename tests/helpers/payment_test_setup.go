package helpers

import (
	"gorm.io/gorm"

	"github.com/reinp/event-platform/backend/internal/database"
	"github.com/reinp/event-platform/backend/internal/service"
)

func SetupPaymentService(
	db *gorm.DB,
) (
	*service.PaymentService,
	*FakePaymentGateway,
) {

	executor := database.NewGormExecutor(
		db,
	)

	purchaseService := NewPurchaseService(
		db,
	)

	ticketService := NewTicketService(
		db,
	)

	fakeGateway := &FakePaymentGateway{}

	paymentService := NewPaymentService(
		executor,
		purchaseService,
		ticketService,
		fakeGateway,
	)

	return paymentService, fakeGateway
}
