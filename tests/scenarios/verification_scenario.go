package scenarios

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	appModels "github.com/reinp/event-platform/backend/internal/models"
	"github.com/reinp/event-platform/backend/internal/models/enum"

	"github.com/reinp/event-platform/backend/tests/fixtures"
	"github.com/reinp/event-platform/backend/tests/helpers"
)

type VerificationScenario struct {
	Staff appModels.User

	StaffMember appModels.PartyMember

	Customer  appModels.User
	OtherUser appModels.User

	Party appModels.Party

	Category appModels.Category

	TicketCategory appModels.TicketCategory

	Window appModels.TicketAccessWindow

	Ticket appModels.Ticket

	Purchase appModels.Purchase

	Scan *appModels.TicketScan
}

func CreateVerificationScenario(
	t *testing.T,
	db *gorm.DB,
	clock *helpers.FakeClock,
	requiresVerification ...bool,
) *VerificationScenario {

	t.Helper()

	verificationRequired := true

	if len(requiresVerification) > 0 {
		verificationRequired = requiresVerification[0]
	}

	// STAFF

	staff := fixtures.User()

	if err := db.Create(&staff).Error; err != nil {
		t.Fatal(err)
	}

	// CUSTOMER

	customer := fixtures.User()

	if err := db.Create(&customer).Error; err != nil {
		t.Fatal(err)
	}

	// OTHER USER

	otherUser := fixtures.User()

	if err := db.Create(&otherUser).Error; err != nil {
		t.Fatal(err)
	}

	// CATEGORY

	category := fixtures.Category()

	if err := db.Create(&category).Error; err != nil {
		t.Fatal(err)
	}

	// PARTY

	party := fixtures.PartyWithOrganizer(
		staff.ID,
	)

	party.CategoryID = category.ID

	if err := db.Create(&party).Error; err != nil {
		t.Fatal(err)
	}

	// PARTY MEMBER (STAFF)

	member := appModels.PartyMember{
		ID: uuid.New(),

		UserID: staff.ID,

		PartyID: party.ID,
	}

	if err := db.Create(&member).Error; err != nil {
		t.Fatal(err)
	}

	// STAFF ROLE

	role := appModels.PartyMemberRole{
		ID: uuid.New(),

		PartyMemberID: member.ID,

		Role: enum.RoleStaff,
	}

	if err := db.Create(&role).Error; err != nil {
		t.Fatal(err)
	}

	// TICKET CATEGORY

	ticketCategory := appModels.TicketCategory{

		ID: uuid.New(),

		Name: "VIP",

		Price: 100,

		Capacity: 100,

		PartyID: party.ID,

		RequiresVerification: verificationRequired,
	}

	if err := db.Create(&ticketCategory).Error; err != nil {
		t.Fatal(err)
	}

	// ACCESS WINDOW

	window := appModels.TicketAccessWindow{

		ID: uuid.New(),

		TicketCategoryID: ticketCategory.ID,

		StartsAt: clock.Current.Add(-time.Hour),

		EndsAt: clock.Current.Add(time.Hour),
	}

	if err := db.Create(&window).Error; err != nil {
		t.Fatal(err)
	}

	// PURCHASE

	purchase := fixtures.Purchase(
		customer.ID,
		party.ID,
	)

	if err := db.Create(&purchase).Error; err != nil {
		t.Fatal(err)
	}

	// TICKET

	ticket := fixtures.Ticket()

	ticket.TicketCategoryID = ticketCategory.ID

	ticket.UserID = customer.ID

	ticket.PurchaseID = purchase.ID

	if err := db.Create(&ticket).Error; err != nil {
		t.Fatal(err)
	}

	return &VerificationScenario{

		Staff: staff,

		StaffMember: member,

		Customer: customer,

		OtherUser: otherUser,

		Party: party,

		Category: category,

		TicketCategory: ticketCategory,

		Window: window,

		Ticket: ticket,
	}
}

func createVerifiedTicketScenario(
	t *testing.T,
	db *gorm.DB,
	clock *helpers.FakeClock,
) *VerificationScenario {

	t.Helper()

	scenario := CreateVerificationScenario(
		t,
		db,
		clock,
		true,
	)

	ticketService := helpers.NewTicketService(
		db,
		clock,
	)

	scan, err := ticketService.Scan(
		context.Background(),
		scenario.Staff.ID,
		scenario.Ticket.Code,
	)

	if err != nil {
		t.Fatal(err)
	}

	if scan.Status != enum.TicketScanPending {
		t.Fatalf(
			"expected pending scan, got %s",
			scan.Status,
		)
	}

	err = ticketService.VerifyScan(
		context.Background(),
		scan.ID,
		scenario.Staff.ID,
		true,
	)

	if err != nil {
		t.Fatal(err)
	}

	scenario.Scan = scan

	return scenario
}

func enableVerification(
	t *testing.T,
	db *gorm.DB,
	category *appModels.TicketCategory,
	value bool,
) {

	if err := db.
		Model(category).
		Update(
			"requires_verification",
			value,
		).
		Error; err != nil {

		t.Fatal(err)
	}
}
