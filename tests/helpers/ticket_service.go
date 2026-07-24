package helpers

import (
	"os"
	"strconv"
	"time"

	"github.com/reinp/event-platform/backend/internal/clock"
	"github.com/reinp/event-platform/backend/internal/database"
	"github.com/reinp/event-platform/backend/internal/repository"
	"github.com/reinp/event-platform/backend/internal/service"
	"gorm.io/gorm"
)

func NewTicketService(
	db *gorm.DB,
	clocks ...clock.Clock,
) *service.TicketService {

	executor := database.NewGormExecutor(db)

	appClock := clock.Clock(
		clock.RealClock{},
	)

	if len(clocks) > 0 {
		appClock = clocks[0]
	}

	return service.NewTicketService(
		repository.NewTicketRepository(executor),
		repository.NewPartyRepository(db),
		repository.NewTicketCategoryRepository(executor),
		repository.NewPartyMemberRepository(executor),
		repository.NewTicketScanRepository(executor),
		repository.NewTicketAccessWindowRepository(executor),
		executor,
		appClock,
		getTicketVerificationTTL(),
	)
}

func getTicketVerificationTTL() time.Duration {

	minutes := 15

	value := os.Getenv(
		"TICKET_VERIFICATION_TTL_MINUTES",
	)

	if value != "" {

		if parsed, err := strconv.Atoi(value); err == nil {
			minutes = parsed
		}
	}

	return time.Duration(minutes) * time.Minute
}
