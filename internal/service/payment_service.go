package service

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

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
	permissionService      *PermissionService
	paymentGateway         payment.Gateway
	paymentEventRepository *repository.PaymentEventRepository
	purchaseRepository     *repository.PurchaseRepository
	refundService          *RefundService
}

func NewPaymentService(
	purchaseService *PurchaseService,
	ticketService *TicketService,
	permissionService *PermissionService,
	paymentGateway payment.Gateway,
	paymentEventRepository *repository.PaymentEventRepository,
	purchaseRepository *repository.PurchaseRepository,
	refundService *RefundService,
) *PaymentService {

	return &PaymentService{
		purchaseService:        purchaseService,
		ticketService:          ticketService,
		permissionService:      permissionService,
		paymentGateway:         paymentGateway,
		paymentEventRepository: paymentEventRepository,
		purchaseRepository:     purchaseRepository,
		refundService:          refundService,
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

	if purchase.Status == enum.PurchaseStatusPaid {
		return "", errors.New(
			"purchase already paid",
		)
	}

	if purchase.PaymentID != "" {
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

	if err := s.purchaseService.AttachPayment(
		ctx,
		purchase.ID,
		"paypal",
		order.ID,
	); err != nil {
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

	if err := s.ticketService.GenerateFromPurchase(
		ctx,
		purchase,
	); err != nil {
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
) (string, error) {

	return s.paymentGateway.CaptureOrder(
		ctx,
		orderID,
	)
}

func (s *PaymentService) RefundPayment(
	ctx context.Context,
	purchaseID uuid.UUID,
	userID uuid.UUID,
) error {

	purchase, err := s.purchaseService.GetPurchase(
		ctx,
		purchaseID,
	)

	if err != nil {

		if errors.Is(
			err,
			gorm.ErrRecordNotFound,
		) {
			return appErrors.ErrPurchaseNotFound
		}

		return err
	}

	if err := s.permissionService.RequireManageRefunds(
		ctx,
		purchase.PartyID,
		userID,
	); err != nil {
		return err
	}

	if purchase.Status == enum.PurchaseStatusRefunded {
		return appErrors.ErrPurchaseAlreadyRefunded
	}

	return s.purchaseRepository.Transaction(
		ctx,
		func(tx database.DBExecutor) error {

			purchase, err := s.purchaseRepository.FindByID(
				tx,
				purchaseID,
			)

			if err != nil {

				if errors.Is(
					err,
					gorm.ErrRecordNotFound,
				) {
					return appErrors.ErrPurchaseNotFound
				}

				return err
			}

			if purchase.Status == enum.PurchaseStatusRefunded {
				return appErrors.ErrPurchaseAlreadyRefunded
			}

			refundAmount, err := s.refundService.CalculateRefundAmount(
				purchase,
			)

			if err != nil {
				return err
			}

			refundID, err := s.paymentGateway.RefundPayment(
				ctx,
				purchase.PaymentID,
				refundAmount,
			)

			if err != nil {
				return err
			}

			now := time.Now()

			purchase.Status = enum.PurchaseStatusRefunded
			purchase.RefundID = refundID
			purchase.RefundProvider = purchase.PaymentProvider
			purchase.RefundedAt = &now

			if err := s.purchaseRepository.Update(
				tx,
				purchase,
			); err != nil {
				return err
			}

			return s.ticketService.CancelByPurchase(
				tx,
				purchase.ID,
			)
		},
	)
}
