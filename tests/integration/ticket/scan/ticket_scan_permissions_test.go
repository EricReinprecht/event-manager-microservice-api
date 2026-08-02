package scan

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/reinp/event-platform/backend/internal/appErrors"
	"github.com/reinp/event-platform/backend/internal/models"
	"github.com/reinp/event-platform/backend/internal/models/enum"

	"github.com/reinp/event-platform/backend/tests/fixtures"
	"github.com/reinp/event-platform/backend/tests/helpers"
	"github.com/reinp/event-platform/backend/tests/scenarios"
)

func TestNonStaffCannotScanTicket(t *testing.T) {

	db, _ := helpers.TestDatabase()
	helpers.CleanDatabase(db)

	clock := helpers.NewFakeClock(time.Now())

	scenario := scenarios.CreateScanScenario(
		t,
		db,
		clock,
	)

	user := fixtures.User()

	if err := db.Create(&user).Error; err != nil {
		t.Fatal(err)
	}

	_, err := scenario.TicketService.Scan(
		context.Background(),
		user.ID,
		scenario.Ticket.Code,
	)

	if !errors.Is(
		err,
		appErrors.ErrNotAllowed,
	) {

		t.Fatalf(
			"expected ErrNotAllowed got %v",
			err,
		)
	}
}

func TestStaffFromDifferentPartyCannotScanTicket(t *testing.T) {

	db, err := helpers.TestDatabase()

	if err != nil {
		t.Fatal(err)
	}

	if err := helpers.CleanDatabase(db); err != nil {
		t.Fatal(err)
	}

	scenario := scenarios.CreateScanScenario(
		t,
		db,
		helpers.NewFakeClock(time.Now()),
	)

	// create external staff user

	otherStaff := fixtures.User()

	if err := db.Create(
		&otherStaff,
	).Error; err != nil {

		t.Fatal(err)
	}

	// create second party where this staff belongs

	otherParty := scenarios.CreateStaffInParty(
		t,
		db,
		otherStaff,
		scenario.Category,
	)

	_ = otherParty

	// make sure ticket belongs to scenario party

	scenario.Ticket.TicketCategoryID =
		scenario.TicketCategory.ID

	if err := db.Save(
		&scenario.Ticket,
	).Error; err != nil {

		t.Fatal(err)
	}

	// staff from another party scans

	_, err = scenario.TicketService.Scan(
		context.Background(),
		otherStaff.ID,
		scenario.Ticket.Code,
	)

	if !errors.Is(
		err,
		appErrors.ErrNotAllowed,
	) {

		t.Fatalf(
			"expected ErrNotAllowed, got %v",
			err,
		)
	}
}

func TestPartyOrganizerCanScanTicket(t *testing.T) {

	db, _ := helpers.TestDatabase()
	helpers.CleanDatabase(db)

	clock := helpers.NewFakeClock(time.Now())

	scenario := scenarios.CreateScanScenario(
		t,
		db,
		clock,
	)

	// organizer membership is created by CreateScanScenario

	scan, err := scenario.TicketService.Scan(
		context.Background(),
		scenario.Organizer.ID,
		scenario.Ticket.Code,
	)

	if err != nil {
		t.Fatal(err)
	}

	if scan.Status != enum.TicketScanVerified {
		t.Fatalf(
			"expected verified scan, got %s",
			scan.Status,
		)
	}

	if scan.ScannedByID != scenario.Organizer.ID {
		t.Fatalf(
			"expected scanner %s, got %s",
			scenario.Organizer.ID,
			scan.ScannedByID,
		)
	}
}

func TestAdminCanScanTicket(t *testing.T) {

	db, err := helpers.TestDatabase()

	if err != nil {
		t.Fatal(err)
	}

	if err := helpers.CleanDatabase(db); err != nil {
		t.Fatal(err)
	}

	scenario := scenarios.CreateScanScenario(
		t,
		db,
		helpers.NewFakeClock(time.Now()),
	)

	// replace staff role with admin role

	scenarios.RemovePartyMember(
		t,
		db,
		&scenario.Member,
	)

	scenario.Member = scenarios.AddPartyRole(
		t,
		db,
		scenario.Staff.ID,
		scenario.Party.ID,
		enum.PartyRoleAdmin,
	)

	// scan

	scan, err := scenario.TicketService.Scan(
		context.Background(),
		scenario.Staff.ID,
		scenario.Ticket.Code,
	)

	if err != nil {
		t.Fatal(err)
	}

	if scan.TicketID != scenario.Ticket.ID {

		t.Fatalf(
			"expected ticket %s, got %s",
			scenario.Ticket.ID,
			scan.TicketID,
		)
	}

	if scan.Status != enum.TicketScanVerified {

		t.Fatalf(
			"expected verified scan, got %s",
			scan.Status,
		)
	}
}

func TestCustomerCannotScanOwnTicket(t *testing.T) {

	db, _ := helpers.TestDatabase()
	helpers.CleanDatabase(db)

	scenario := scenarios.CreateScanScenario(
		t,
		db,
		helpers.NewFakeClock(time.Now()),
	)

	_, err := scenario.TicketService.Scan(
		context.Background(),
		scenario.Customer.ID,
		scenario.Ticket.Code,
	)

	if !errors.Is(
		err,
		appErrors.ErrNotAllowed,
	) {

		t.Fatalf(
			"expected ErrNotAllowed got %v",
			err,
		)
	}

	var count int64

	db.Model(
		&models.TicketScan{},
	).
		Where(
			"ticket_id = ?",
			scenario.Ticket.ID,
		).
		Count(&count)

	if count != 0 {
		t.Fatalf(
			"expected no scans got %d",
			count,
		)
	}
}

func TestScanWithNoPartyMemberRecord(t *testing.T) {

	db, err := helpers.TestDatabase()

	if err != nil {
		t.Fatal(err)
	}

	if err := helpers.CleanDatabase(db); err != nil {
		t.Fatal(err)
	}

	clock := helpers.NewFakeClock(time.Now().UTC())

	scenario := scenarios.CreateScanScenario(
		t,
		db,
		clock,
	)

	// REMOVE ALL MEMBERSHIP FOR SCANNER

	var member models.PartyMember

	if err := db.
		Where(
			"user_id = ? AND party_id = ?",
			scenario.Staff.ID,
			scenario.Party.ID,
		).
		First(&member).
		Error; err != nil {

		t.Fatal(err)
	}

	// remove roles first
	if err := db.
		Where(
			"party_member_id = ?",
			member.ID,
		).
		Delete(
			&models.PartyMemberRole{},
		).
		Error; err != nil {

		t.Fatal(err)
	}

	// now remove membership
	if err := db.
		Delete(
			&member,
		).
		Error; err != nil {

		t.Fatal(err)
	}

	// TRY SCAN WITHOUT MEMBERSHIP

	_, err = scenario.TicketService.Scan(
		context.Background(),
		scenario.Staff.ID,
		scenario.Ticket.Code,
	)

	if !errors.Is(
		err,
		appErrors.ErrNotAllowed,
	) {

		t.Fatalf(
			"expected ErrNotAllowed, got %v",
			err,
		)
	}
}

func TestStaffCanScanTicketOwnedByAnotherUser(t *testing.T) {

	db, err := helpers.TestDatabase()

	if err != nil {
		t.Fatal(err)
	}

	if err := helpers.CleanDatabase(db); err != nil {
		t.Fatal(err)
	}

	scenario := scenarios.CreateScanScenario(
		t,
		db,
		helpers.NewFakeClock(time.Now()),
	)

	// Make sure ticket belongs to another user
	scenario.Ticket.UserID = scenario.Customer.ID

	if err := db.Save(
		&scenario.Ticket,
	).Error; err != nil {

		t.Fatal(err)
	}

	// STAFF SCANS OTHER USER'S TICKET

	scan, err := scenario.TicketService.Scan(
		context.Background(),
		scenario.Staff.ID,
		scenario.Ticket.Code,
	)

	if err != nil {
		t.Fatal(err)
	}

	// VERIFY

	if scan.TicketID != scenario.Ticket.ID {

		t.Fatalf(
			"expected ticket %s got %s",
			scenario.Ticket.ID,
			scan.TicketID,
		)
	}

	if scan.Status != enum.TicketScanVerified {

		t.Fatalf(
			"expected verified scan, got %s",
			scan.Status,
		)
	}

	if scan.ScannedByID != scenario.Staff.ID {

		t.Fatalf(
			"expected scanner %s got %s",
			scenario.Staff.ID,
			scan.ScannedByID,
		)
	}

	// OWNER MUST NOT BE SCANNER

	if scenario.Ticket.UserID == scenario.Staff.ID {

		t.Fatal(
			"test setup invalid: ticket owner equals scanner",
		)
	}
}
