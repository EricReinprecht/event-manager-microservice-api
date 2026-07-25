package verification

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/reinp/event-platform/backend/internal/appErrors"
	appModels "github.com/reinp/event-platform/backend/internal/models"
	"github.com/reinp/event-platform/backend/internal/models/enum"

	"github.com/reinp/event-platform/backend/tests/helpers"
	"github.com/reinp/event-platform/backend/tests/scenarios"
)

func TestConcurrentVerifyAndReject(t *testing.T) {

	db, err := helpers.TestDatabase()

	if err != nil {
		t.Fatal(err)
	}

	if err := helpers.CleanDatabase(db); err != nil {
		t.Fatal(err)
	}

	clock := &helpers.FakeClock{
		Current: time.Now().
			UTC().
			Truncate(time.Microsecond),
	}

	ticketService := helpers.NewTicketService(
		db,
		clock,
	)

	scenario := scenarios.CreateVerificationScenario(
		t,
		db,
		clock,
		true,
	)

	staff := scenario.Staff
	ticket := scenario.Ticket

	// CREATE PENDING SCAN

	scan, err := ticketService.Scan(
		context.Background(),
		staff.ID,
		ticket.Code,
	)

	if err != nil {
		t.Fatal(err)
	}

	scenarios.AssertScanStatus(
		t,
		scan,
		enum.TicketScanPending,
	)

	// CONCURRENT DECISIONS

	var wg sync.WaitGroup

	results := make(chan error, 2)

	wg.Add(2)

	go func() {

		defer wg.Done()

		results <- ticketService.VerifyScan(
			context.Background(),
			scan.ID,
			staff.ID,
			true,
		)

	}()

	go func() {

		defer wg.Done()

		results <- ticketService.VerifyScan(
			context.Background(),
			scan.ID,
			staff.ID,
			false,
		)

	}()

	wg.Wait()

	close(results)

	success := 0
	failures := 0

	for err := range results {

		if err == nil {
			success++
		} else {
			failures++
		}
	}

	if success != 1 {

		t.Fatalf(
			"expected exactly one successful decision, got %d",
			success,
		)
	}

	if failures != 1 {

		t.Fatalf(
			"expected exactly one failed decision, got %d",
			failures,
		)
	}

	// FINAL STATE

	var updated appModels.TicketScan

	if err := db.
		First(
			&updated,
			"id = ?",
			scan.ID,
		).
		Error; err != nil {

		t.Fatal(err)
	}

	switch updated.Status {

	case enum.TicketScanVerified:

		if updated.VerifiedAt == nil {
			t.Fatal(
				"verified scan missing decision timestamp",
			)
		}

		if updated.VerifiedByID == nil {
			t.Fatal(
				"verified scan missing decision user",
			)
		}

		if updated.VerifiedUntil == nil {
			t.Fatal(
				"verified scan missing verification expiry",
			)
		}

	case enum.TicketScanRejected:

		if updated.VerifiedAt == nil {
			t.Fatal(
				"rejected scan missing decision timestamp",
			)
		}

		if updated.VerifiedByID == nil {
			t.Fatal(
				"rejected scan missing rejecting user",
			)
		}

		if updated.VerifiedUntil != nil {
			t.Fatal(
				"rejected scan should not have verification validity",
			)
		}

	default:

		t.Fatalf(
			"unexpected final status %s",
			updated.Status,
		)
	}
}

func TestConcurrentVerifySameScan(t *testing.T) {

	db, err := helpers.TestDatabase()

	if err != nil {
		t.Fatal(err)
	}

	if err := helpers.CleanDatabase(db); err != nil {
		t.Fatal(err)
	}

	clock := &helpers.FakeClock{
		Current: time.Now().
			UTC().
			Truncate(time.Microsecond),
	}

	ticketService := helpers.NewTicketService(
		db,
		clock,
	)

	scenario := scenarios.CreateVerificationScenario(
		t,
		db,
		clock,
		true,
	)

	staff1 := scenario.Staff

	// ADD SECOND STAFF MEMBER

	staff2 := scenarios.AddSecondPartyStaff(
		t,
		db,
		scenario.Party.ID,
	)

	ticket := scenario.Ticket

	// CREATE PENDING SCAN

	scan, err := ticketService.Scan(
		context.Background(),
		staff1.ID,
		ticket.Code,
	)

	if err != nil {
		t.Fatal(err)
	}

	scenarios.AssertScanStatus(
		t,
		scan,
		enum.TicketScanPending,
	)

	// RUN TWO APPROVALS CONCURRENTLY

	results := make(chan error, 2)

	var wg sync.WaitGroup

	wg.Add(2)

	go func() {

		defer wg.Done()

		results <- ticketService.VerifyScan(
			context.Background(),
			scan.ID,
			staff1.ID,
			true,
		)

	}()

	go func() {

		defer wg.Done()

		results <- ticketService.VerifyScan(
			context.Background(),
			scan.ID,
			staff2.ID,
			true,
		)

	}()

	wg.Wait()

	close(results)

	successes := 0
	failures := 0

	for err := range results {

		if err == nil {

			successes++

		} else if errors.Is(
			err,
			appErrors.ErrTicketScanAlreadyDecided,
		) {

			failures++

		} else {

			t.Fatalf(
				"unexpected error: %v",
				err,
			)
		}
	}

	if successes != 1 {

		t.Fatalf(
			"expected exactly one successful verification, got %d",
			successes,
		)
	}

	if failures != 1 {

		t.Fatalf(
			"expected exactly one rejected verification attempt, got %d",
			failures,
		)
	}

	// VERIFY FINAL STATE

	var updated appModels.TicketScan

	if err := db.First(
		&updated,
		"id = ?",
		scan.ID,
	).Error; err != nil {

		t.Fatal(err)
	}

	scenarios.AssertScanStatus(
		t,
		&updated,
		enum.TicketScanVerified,
	)

	if updated.VerifiedByID == nil {

		t.Fatal(
			"expected verifying staff metadata",
		)
	}
}
