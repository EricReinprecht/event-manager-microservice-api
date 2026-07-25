package service

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"

	"github.com/reinp/event-platform/backend/internal/appErrors"
	"github.com/reinp/event-platform/backend/internal/database"
	"github.com/reinp/event-platform/backend/internal/models"
	"github.com/reinp/event-platform/backend/internal/models/enum"
	"github.com/reinp/event-platform/backend/internal/repository"
	"github.com/reinp/event-platform/backend/internal/requests"
)

type PurchaseService struct {
	repository       *repository.PurchaseRepository
	ticketRepository *repository.TicketRepository
}

func NewPurchaseService(
	repository *repository.PurchaseRepository,
	ticketRepository *repository.TicketRepository,
) *PurchaseService {

	return &PurchaseService{
		repository:       repository,
		ticketRepository: ticketRepository,
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
				Status:  enum.PruchaseStatusPending,

				ExpiresAt: time.Now().Add(
					30 * time.Minute,
				),
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

				// -----------------------------
				// Sold-out protection
				// -----------------------------

				sold, err := s.ticketRepository.CountByCategoryTx(
					tx,
					category.ID,
				)

				if err != nil {
					return err
				}

				if sold+int64(item.Quantity) > int64(category.Capacity) {

					return appErrors.ErrTicketSoldOut
				}

				// -----------------------------
				// Create purchase item snapshot
				// -----------------------------

				purchase.Items = append(
					purchase.Items,
					models.PurchaseItem{
						ID: uuid.New(),

						TicketCategoryID: category.ID,

						Quantity: item.Quantity,

						// price snapshot
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

			if purchase.Status != enum.PruchaseStatusPending {
				return errors.New("purchase is not pending")
			}

			if time.Now().After(
				purchase.ExpiresAt,
			) {

				return errors.New(
					"purchase expired",
				)
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

func (s *PurchaseService) GetPurchase(
	ctx context.Context,
	id uuid.UUID,
) (*models.Purchase, error) {

	return s.repository.Find(
		ctx,
		id,
	)
}

func (s *PurchaseService) FindByPaymentID(
	ctx context.Context,
	paymentID string,
) (*models.Purchase, error) {

	return s.repository.FindByPaymentID(
		ctx,
		paymentID,
	)
}

func (s *PurchaseService) ConfirmPayment(
	ctx context.Context,
	paymentID string,
) (*models.Purchase, error) {

	var result *models.Purchase

	err := s.repository.Transaction(
		ctx,
		func(tx database.DBExecutor) error {

			purchase, err := s.repository.FindByPaymentIDForUpdate(
				tx,
				paymentID,
			)

			if err != nil {
				return err
			}

			if purchase.Status == enum.PurchaseStatusPaid {

				result = purchase
				return nil
			}

			purchase.Status = enum.PurchaseStatusPaid

			if err := s.repository.Update(
				tx,
				purchase,
			); err != nil {

				return err
			}

			result = purchase

			return nil
		},
	)

	if err != nil {
		return nil, err
	}

	return result, nil
}
