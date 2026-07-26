package service

import (
	"errors"

	"github.com/reinp/event-platform/backend/internal/models"
)

type RefundService struct {
}

func NewRefundService() *RefundService {
	return &RefundService{}
}

// CalculateRefundAmount calculates how much money should be refunded
// based on the refund policy of each ticket category.
func (s *RefundService) CalculateRefundAmount(
	purchase *models.Purchase,
) (int64, error) {

	var amount int64

	for _, item := range purchase.Items {

		if item.TicketCategory.RefundPolicy == nil {
			continue
		}

		percentage := item.TicketCategory.RefundPolicy.Percentage

		amount +=
			item.UnitPrice *
				int64(item.Quantity) *
				int64(percentage) /
				100
	}

	if amount <= 0 {
		return 0, errors.New("no refundable amount")
	}

	return amount, nil
}
