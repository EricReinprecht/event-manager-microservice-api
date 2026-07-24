package service

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/reinp/event-platform/backend/internal/appErrors"
	"github.com/reinp/event-platform/backend/internal/models"
	"github.com/reinp/event-platform/backend/internal/models/enum"
	"github.com/reinp/event-platform/backend/internal/requests"
)

type PurchaseService struct {
	db *gorm.DB
}

func NewPurchaseService(
	db *gorm.DB,
) *PurchaseService {

	return &PurchaseService{
		db: db,
	}
}

func (s *PurchaseService) CreatePurchase(
	ctx context.Context,
	userID uuid.UUID,
	partyID uuid.UUID,
	items []requests.PurchaseItemRequest,
) (*models.Purchase, error) {

	var purchase models.Purchase

	err := s.db.Transaction(func(tx *gorm.DB) error {

		purchase = models.Purchase{
			ID:      uuid.New(),
			UserID:  userID,
			PartyID: partyID,
			Status:  enum.StatusPending,
		}

		var total int64

		for _, item := range items {

			var category models.TicketCategory

			if err := tx.First(
				&category,
				"id = ?",
				item.TicketCategoryID,
			).Error; err != nil {

				return appErrors.ErrTicketCategoryNotFound
			}

			if category.PartyID != partyID {
				return appErrors.ErrTicketCategoryNotFound
			}

			price := category.Price

			purchase.Items = append(
				purchase.Items,
				models.PurchaseItem{
					ID:               uuid.New(),
					TicketCategoryID: category.ID,
					Quantity:         item.Quantity,
					UnitPrice:        price,
				},
			)

			total += price * int64(item.Quantity)
		}

		purchase.TotalPrice = total

		return tx.Create(&purchase).Error
	})

	if err != nil {
		return nil, err
	}

	return &purchase, nil
}
