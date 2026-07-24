package service

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/reinp/event-platform/backend/internal/models"
	"github.com/reinp/event-platform/backend/internal/models/enum"
	"github.com/reinp/event-platform/backend/internal/payment/paypal"
)

type PaymentService struct {
	purchaseService *PurchaseService
	ticketService   *TicketService
	paypalClient    *paypal.Client
}

func NewPaymentService(
	purchaseService *PurchaseService,
	ticketService *TicketService,
	paypalClient *paypal.Client,
) *PaymentService {

	return &PaymentService{
		purchaseService: purchaseService,
		ticketService:   ticketService,
		paypalClient:    paypalClient,
	}
}
func (s *PaymentService) CreateCheckout(
	ctx context.Context,
	purchaseID uuid.UUID,
) (string, error) {

	purchase, err := s.purchaseService.GetPurchase(
		ctx,
		purchaseID,
	)

	if err != nil {
		return "", err
	}

	if purchase.Status != enum.StatusPending {
		return "", errors.New(
			"purchase is not pending",
		)
	}

	order, err := s.paypalClient.CreateOrder(
		ctx,
		purchase.TotalPrice,
	)

	if err != nil {
		return "", err
	}

	err = s.purchaseService.AttachPayment(
		ctx,
		purchase.ID,
		"paypal",
		order.ID,
	)

	if err != nil {
		return "", err
	}

	return order.ApprovalURL, nil
}

func (s *PaymentService) ConfirmPayment(
	ctx context.Context,
	paymentID string,
) (*models.Purchase, error) {

	purchase, err := s.purchaseService.ConfirmPayment(
		ctx,
		paymentID,
	)

	if err != nil {
		return nil, err
	}

	err = s.ticketService.GenerateFromPurchase(
		ctx,
		purchase,
	)

	if err != nil {
		return nil, err
	}

	return purchase, nil
}
