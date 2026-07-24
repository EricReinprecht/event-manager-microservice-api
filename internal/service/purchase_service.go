package service

import (
	"context"
	"errors"

	"github.com/google/uuid"

	"github.com/reinp/event-platform/backend/internal/appErrors"
	"github.com/reinp/event-platform/backend/internal/database"
	"github.com/reinp/event-platform/backend/internal/models"
	"github.com/reinp/event-platform/backend/internal/models/enum"
	"github.com/reinp/event-platform/backend/internal/repository"
	"github.com/reinp/event-platform/backend/internal/requests"
)

type PurchaseService struct {
	repository *repository.PurchaseRepository
}

func NewPurchaseService(
	repository *repository.PurchaseRepository,
) *PurchaseService {

	return &PurchaseService{
		repository: repository,
	}
}

func (s *PurchaseService) CreatePurchase(
	ctx context.Context,
	userID uuid.UUID,
	partyID uuid.UUID,
	items []requests.PurchaseItemRequest,
) (*models.Purchase, error) {

	var purchase models.Purchase

	err := s.repository.Transaction(
		ctx,
		func(tx database.DBExecutor) error {

			purchase = models.Purchase{
				ID:      uuid.New(),
				UserID:  userID,
				PartyID: partyID,
				Status:  enum.StatusPending,
			}

			var total int64

			for _, item := range items {

				category, err := s.repository.FindTicketCategory(
					tx,
					item.TicketCategoryID,
				)

				if err != nil {
					return appErrors.ErrTicketCategoryNotFound
				}

				if category.PartyID != partyID {
					return appErrors.ErrTicketCategoryNotFound
				}

				purchase.Items = append(
					purchase.Items,
					models.PurchaseItem{
						ID: uuid.New(),

						TicketCategoryID: category.ID,

						Quantity: item.Quantity,

						UnitPrice: category.Price,
					},
				)

				total += category.Price * int64(item.Quantity)
			}

			purchase.TotalPrice = total

			return s.repository.Create(
				tx,
				&purchase,
			)
		},
	)

	if err != nil {
		return nil, err
	}

	return &purchase, nil
}

func (s *PurchaseService) AttachPayment(
	ctx context.Context,
	purchaseID uuid.UUID,
	provider string,
	paymentID string,
) error {

	return s.repository.Transaction(
		ctx,
		func(tx database.DBExecutor) error {

			purchase, err := s.repository.FindByID(
				tx,
				purchaseID,
			)

			if err != nil {
				return err
			}

			if purchase.Status != enum.StatusPending {
				return errors.New("purchase is not pending")
			}

			return s.repository.UpdatePayment(
				tx,
				purchase,
				provider,
				paymentID,
			)
		},
	)
}

func (s *PurchaseService) ConfirmPayment(
	ctx context.Context,
	paymentID string,
) (*models.Purchase, error) {

	var purchase *models.Purchase

	err := s.repository.Transaction(
		ctx,
		func(tx database.DBExecutor) error {

			var err error

			purchase, err = s.repository.FindByPaymentID(
				tx,
				paymentID,
			)

			if err != nil {
				return err
			}

			// webhook idempotency
			if purchase.Status == enum.StatusPaid {
				return nil
			}

			if purchase.Status != enum.StatusPending {
				return errors.New(
					"purchase cannot be paid",
				)
			}

			return s.repository.UpdateStatus(
				tx,
				purchase,
				enum.StatusPaid,
			)
		},
	)

	if err != nil {
		return nil, err
	}

	return purchase, nil
}

func (s *PurchaseService) GetPurchase(
	ctx context.Context,
	id uuid.UUID,
) (*models.Purchase, error) {

	return s.repository.Find(
		ctx,
		id,
	)
}
