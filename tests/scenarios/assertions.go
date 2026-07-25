package scenarios

import (
	"testing"

	"github.com/google/uuid"
	appModels "github.com/reinp/event-platform/backend/internal/models"
	"github.com/reinp/event-platform/backend/internal/models/enum"
)

func AssertScanStatus(
	t *testing.T,
	scan *appModels.TicketScan,
	expected enum.TicketScanStatus,
) {

	if scan.Status != expected {

		t.Fatalf(
			"expected %s got %s",
			expected,
			scan.Status,
		)
	}
}

func AssertVerified(
	t *testing.T,
	scan *appModels.TicketScan,
	userID uuid.UUID,
) {

	if scan.Status != enum.TicketScanVerified {
		t.Fatalf(
			"expected verified, got %s",
			scan.Status,
		)
	}

	if scan.VerifiedByID == nil {
		t.Fatal("expected verifier")
	}

	if *scan.VerifiedByID != userID {
		t.Fatal("wrong verifier")
	}

	if scan.VerifiedAt == nil {
		t.Fatal("expected verification timestamp")
	}
}

func AssertRejected(
	t *testing.T,
	scan *appModels.TicketScan,
	userID uuid.UUID,
) {

	if scan.Status != enum.TicketScanRejected {
		t.Fatalf(
			"expected rejected, got %s",
			scan.Status,
		)
	}

	if scan.VerifiedByID == nil {
		t.Fatal("expected rejecting user")
	}

	if *scan.VerifiedByID != userID {
		t.Fatal("wrong rejecting user")
	}

	if scan.VerifiedAt == nil {
		t.Fatal("expected rejection timestamp")
	}
}
