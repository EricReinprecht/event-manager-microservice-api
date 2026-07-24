package service

import (
	"context"

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
