package scenarios

import (
	"context"
	"testing"

	"github.com/google/uuid"
	appModels "github.com/reinp/event-platform/backend/internal/models"
)

func CreatePendingScan(
	t *testing.T,
	service interface {
		Scan(context.Context, uuid.UUID, string) (*appModels.TicketScan, error)
	},
	userID uuid.UUID,
	ticket appModels.Ticket,
) *appModels.TicketScan {

	scan, err := service.Scan(
		context.Background(),
		userID,
		ticket.Code,
	)

	if err != nil {
		t.Fatal(err)
	}

	return scan
}

func VerifyScan(
	t *testing.T,
	service interface {
		VerifyScan(context.Context, uuid.UUID, uuid.UUID, bool) error
	},
	scan *appModels.TicketScan,
	userID uuid.UUID,
	approve bool,
) {

	err := service.VerifyScan(
		context.Background(),
		scan.ID,
		userID,
		approve,
	)

	if err != nil {
		t.Fatal(err)
	}
}
