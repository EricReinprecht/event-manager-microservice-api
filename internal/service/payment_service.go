package service

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	appErrors "github.com/reinp/event-platform/backend/internal/appErrors"
	"github.com/reinp/event-platform/backend/internal/database"
	"github.com/reinp/event-platform/backend/internal/models"
	"github.com/reinp/event-platform/backend/internal/models/enum"
	"github.com/reinp/event-platform/backend/internal/payment"
	"github.com/reinp/event-platform/backend/internal/payment/paypal"
	"github.com/reinp/event-platform/backend/internal/repository"
)

type PaymentService struct {
	purchaseService        *PurchaseService
	ticketService          *TicketService
	partyMemberService     *PartyMemberService
	paymentGateway         payment.Gateway
	paymentEventRepository *repository.PaymentEventRepository
	purchaseRepository     *repository.PurchaseRepository
}

func NewPaymentService(
	purchaseService *PurchaseService,
	ticketService *TicketService,
	partyMemberService *PartyMemberService,
	paymentGateway payment.Gateway,
	paymentEventRepository *repository.PaymentEventRepository,
	purchaseRepository *repository.PurchaseRepository,
) *PaymentService {

	return &PaymentService{
		purchaseService:        purchaseService,
		ticketService:          ticketService,
		partyMemberService:     partyMemberService,
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
	if purchase.Status == enum.PurchaseStatusPaid {
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

	var total int64

	for _, item := range purchase.Items {

		total += item.UnitPrice * int64(item.Quantity)

	}

	if total <= 0 {
		return "", errors.New(
			"purchase total must be greater than zero",
		)
	}

	order, err := s.paymentGateway.CreateOrder(
		ctx,
		total,
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
	if purchase.Status == enum.PurchaseStatusPaid {

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

func (s *PaymentService) RefundPayment(
	ctx context.Context,
	purchaseID uuid.UUID,
) error {

	return s.purchaseRepository.Transaction(
		ctx,
		func(tx database.DBExecutor) error {

			purchase, err := s.purchaseRepository.FindByID(
				tx,
				purchaseID,
			)

			if err != nil {
				return err
			}

			if purchase.Status == enum.PurchaseStatusRefunded {
				return errors.New(
					"purchase already refunded",
				)
			}

			refundID, err := s.paymentGateway.RefundPayment(
				ctx,
				purchase.PaymentID,
			)

			if err != nil {
				return err
			}

			now := time.Now()

			purchase.Status = enum.PurchaseStatusRefunded

			purchase.RefundID = refundID

			purchase.RefundProvider = purchase.PaymentProvider

			purchase.RefundedAt = &now

			err = s.purchaseRepository.Update(
				tx,
				purchase,
			)

			if err != nil {
				return err
			}

			return s.ticketService.CancelByPurchase(
				tx,
				purchase.ID,
			)
		},
	)
}
