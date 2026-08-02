package service

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"

	appErrors "github.com/reinp/event-platform/backend/internal/appErrors"
	"github.com/reinp/event-platform/backend/internal/clock"
	"github.com/reinp/event-platform/backend/internal/dto"
	"github.com/reinp/event-platform/backend/internal/models"
	"github.com/reinp/event-platform/backend/internal/models/enum"
	"github.com/reinp/event-platform/backend/internal/repository"
)

type PurchaseService struct {
	repository  *repository.PurchaseRepository
	unitOfWork  *repository.PurchaseUnitOfWork
	clock       clock.Clock
	purchaseTTL time.Duration
}

func NewPurchaseService(
	repository *repository.PurchaseRepository,
	unitOfWork *repository.PurchaseUnitOfWork,
	clock clock.Clock,
	purchaseTTL time.Duration,
) *PurchaseService {

	return &PurchaseService{
		repository:  repository,
		unitOfWork:  unitOfWork,
		clock:       clock,
		purchaseTTL: purchaseTTL,
	}
}

func (s *PurchaseService) CreatePurchase(
	ctx context.Context,
	userID uuid.UUID,
	partyID uuid.UUID,
	items []dto.PurchaseItemRequest,
) (*models.Purchase, error) {

	if len(items) == 0 {
		return nil, appErrors.ErrPurchaseItemsRequired
	}

	purchase := &models.Purchase{
		ID:      uuid.New(),
		UserID:  userID,
		PartyID: partyID,
		Status:  enum.PurchaseStatusPending,

		ExpiresAt: s.clock.Now().Add(
			s.purchaseTTL,
		),

		Items: make(
			[]models.PurchaseItem,
			0,
			len(items),
		),
	}

	err := s.unitOfWork.Transaction(
		ctx,
		func(
			repositories *repository.PurchaseTransactionRepositories,
		) error {

			var total int64

			for _, item := range items {

				category, err :=
					repositories.TicketCategories.FindByID(
						ctx,
						item.TicketCategoryID,
					)

				if err != nil {
					return appErrors.ErrTicketCategoryNotFound
				}

				if category.PartyID != partyID {
					return appErrors.ErrTicketCategoryNotFound
				}

				sold, err :=
					repositories.Tickets.CountByCategory(
						ctx,
						category.ID,
					)

				if err != nil {
					return err
				}

				if sold+int64(item.Quantity) >
					int64(category.Capacity) {

					return appErrors.ErrTicketSoldOut
				}

				purchaseItem := models.PurchaseItem{
					ID: uuid.New(),

					PurchaseID: purchase.ID,

					TicketCategoryID: category.ID,

					Quantity: item.Quantity,

					UnitPrice: category.Price,
				}

				purchase.Items = append(
					purchase.Items,
					purchaseItem,
				)

				total +=
					category.Price *
						int64(item.Quantity)
			}

			purchase.TotalPrice = total

			return repositories.Purchases.Create(
				ctx,
				purchase,
			)
		},
	)

	if err != nil {
		return nil, err
	}

	return purchase, nil
}

func (s *PurchaseService) AttachPayment(
	ctx context.Context,
	purchaseID uuid.UUID,
	provider string,
	paymentID string,
) error {

	return s.unitOfWork.Transaction(
		ctx,
		func(
			repositories *repository.PurchaseTransactionRepositories,
		) error {

			purchase, err :=
				repositories.Purchases.FindByIDForUpdate(
					ctx,
					purchaseID,
				)

			if err != nil {
				return err
			}

			if purchase.Status !=
				enum.PurchaseStatusPending {

				return appErrors.ErrPurchaseNotPending
			}

			if s.clock.Now().After(
				purchase.ExpiresAt,
			) {
				return appErrors.ErrPurchaseExpired
			}

			if purchase.PaymentID != "" {
				return appErrors.ErrCheckoutAlreadyCreated
			}

			return repositories.Purchases.UpdatePayment(
				ctx,
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

	return s.repository.FindByID(
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

	err := s.unitOfWork.Transaction(
		ctx,
		func(
			repositories *repository.PurchaseTransactionRepositories,
		) error {

			purchase, err :=
				repositories.Purchases.
					FindByPaymentIDForUpdate(
						ctx,
						paymentID,
					)

			if err != nil {

				if errors.Is(
					err,
					appErrors.ErrPurchaseNotFound,
				) {
					return appErrors.ErrUnknownPaymentOrder
				}

				return err
			}

			if purchase.Status ==
				enum.PurchaseStatusPaid {

				result = purchase

				return nil
			}

			if purchase.Status !=
				enum.PurchaseStatusPending {

				return appErrors.ErrPurchaseNotPending
			}

			purchase.Status =
				enum.PurchaseStatusPaid

			if err := repositories.Purchases.Update(
				ctx,
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
