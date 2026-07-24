package service

import (
	"context"
	"errors"
	"log"

	"github.com/google/uuid"
	"github.com/reinp/event-platform/backend/internal/models"
	"github.com/reinp/event-platform/backend/internal/models/enum"
	"github.com/reinp/event-platform/backend/internal/payment"
	"github.com/reinp/event-platform/backend/internal/payment/paypal"
	"github.com/reinp/event-platform/backend/internal/repository"
)

type PaymentService struct {
	purchaseService        *PurchaseService
	ticketService          *TicketService
	paymentGateway         payment.Gateway
	paymentEventRepository *repository.PaymentEventRepository
}

func NewPaymentService(
	purchaseService *PurchaseService,
	ticketService *TicketService,
	paymentGateway payment.Gateway,
	paymentEventRepository *repository.PaymentEventRepository,
) *PaymentService {

	return &PaymentService{
		purchaseService:        purchaseService,
		ticketService:          ticketService,
		paymentGateway:         paymentGateway,
		paymentEventRepository: paymentEventRepository,
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

	// Already paid
	if purchase.Status == enum.StatusPaid {
		return "", errors.New(
			"purchase already paid",
		)
	}

	// Already has PayPal order
	if purchase.PaymentID != "" {

		// optionally return existing approval URL later
		return "", errors.New(
			"checkout already created",
		)
	}

	order, err := s.paymentGateway.CreateOrder(
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

	log.Println(
		"ConfirmPayment called with:",
		paymentID,
	)

	purchase, err := s.purchaseService.FindByPaymentID(
		ctx,
		paymentID,
	)

	if err != nil {
		log.Println(
			"PAYMENT NOT FOUND:",
			paymentID,
		)

		return nil, err
	}

	log.Println(
		"purchase found:",
		purchase.ID,
		purchase.Status,
	)

	// Idempotency protection
	if purchase.Status == enum.StatusPaid {

		log.Println(
			"purchase already paid",
		)

		return purchase, nil
	}

	// Payment was already captured by PayPal webhook.
	// Only update local state now.

	purchase, err = s.purchaseService.ConfirmPayment(
		ctx,
		paymentID,
	)

	if err != nil {
		return nil, err
	}

	log.Println(
		"purchase marked paid",
	)

	err = s.ticketService.GenerateFromPurchase(
		ctx,
		purchase,
	)

	if err != nil {
		return nil, err
	}

	log.Println(
		"tickets generated",
	)

	return purchase, nil
}

func (s *PaymentService) VerifyWebhook(
	ctx context.Context,
	headers paypal.WebhookHeaders,
	body []byte,
) error {

	return s.paymentGateway.VerifyWebhookSignature(
		ctx,
		headers,
		body,
	)
}

func (s *PaymentService) FindPaymentEvent(
	ctx context.Context,
	eventID string,
) (*models.PaymentEvent, error) {

	return s.paymentEventRepository.FindByEventID(
		ctx,
		eventID,
	)
}

func (s *PaymentService) CreatePaymentEvent(
	ctx context.Context,
	event *models.PaymentEvent,
) error {

	return s.paymentEventRepository.Create(
		event,
	)
}

func (s *PaymentService) MarkPaymentEventProcessed(
	ctx context.Context,
	eventID string,
) error {

	event, err := s.paymentEventRepository.FindByEventID(
		ctx,
		eventID,
	)

	if err != nil {
		return err
	}

	event.Processed = true

	return s.paymentEventRepository.Update(
		event,
	)
}

func (s *PaymentService) CapturePayment(
	ctx context.Context,
	orderID string,
) error {

	return s.paymentGateway.CaptureOrder(
		ctx,
		orderID,
	)
}
