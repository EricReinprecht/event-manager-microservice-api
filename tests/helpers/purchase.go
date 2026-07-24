package helpers

import (
	"gorm.io/gorm"

	"github.com/reinp/event-platform/backend/internal/service"
)

func NewPurchaseService(
	db *gorm.DB,
) *service.PurchaseService {

	return service.NewPurchaseService(
		db,
	)
}
