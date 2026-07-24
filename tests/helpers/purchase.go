package helpers

import (
	"gorm.io/gorm"

	"github.com/reinp/event-platform/backend/internal/database"
	"github.com/reinp/event-platform/backend/internal/repository"
	"github.com/reinp/event-platform/backend/internal/service"
)

func NewPurchaseService(
	db *gorm.DB,
) *service.PurchaseService {

	executor := database.NewGormExecutor(db)

	repository := repository.NewPurchaseRepository(
		executor,
	)

	return service.NewPurchaseService(
		repository,
	)
}
