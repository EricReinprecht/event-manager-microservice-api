package service

import (
	"context"
	"errors"

	"github.com/google/uuid"
	appErrors "github.com/reinp/event-platform/backend/internal/appErrors"
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
	purchaseRepository     *repository.PurchaseRepository
}

func NewPaymentService(
	purchaseService *PurchaseService,
	ticketService *TicketService,
	paymentGateway payment.Gateway,
	paymentEventRepository *repository.PaymentEventRepository,
	purchaseRepository *repository.PurchaseRepository,
) *PaymentService {

	return &PaymentService{
		purchaseService:        purchaseService,
		ticketService:          ticketService,
		paymentGateway:         paymentGateway,
		paymentEventRepository: paymentEventRepository,
		purchaseRepository:     purchaseRepository,
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

	purchase, err := s.purchaseService.FindByPaymentID(
		ctx,
		paymentID,
	)

	if err != nil {

		if errors.Is(
			err,
			appErrors.ErrPurchaseNotFound,
		) {
			return nil, appErrors.ErrUnknownPaymentOrder
		}

		return nil, err
	}

	// idempotency
	if purchase.Status == enum.StatusPaid {

		return purchase, nil
	}

	purchase, err = s.purchaseService.ConfirmPayment(
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
