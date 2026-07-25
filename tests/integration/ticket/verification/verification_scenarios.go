package verification

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
}

func createVerificationScenario(
	t *testing.T,
	db *gorm.DB,
	clock *helpers.FakeClock,
) *VerificationScenario {

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

		RequiresVerification: true,
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

func changeStaffRole(
	t *testing.T,
	db *gorm.DB,
	s *VerificationScenario,
	role enum.PartyRole,
) {

	if err := db.
		Model(&appModels.PartyMemberRole{}).
		Where(
			"party_member_id = ?",
			s.StaffMember.ID,
		).
		Update(
			"role",
			role,
		).
		Error; err != nil {

		t.Fatal(err)
	}
}

func makeOrganizer(
	t *testing.T,
	db *gorm.DB,
	s *VerificationScenario,
) {
	changeStaffRole(
		t,
		db,
		s,
		enum.RoleOrganizer,
	)
}

func makeAdmin(
	t *testing.T,
	db *gorm.DB,
	s *VerificationScenario,
) {
	changeStaffRole(
		t,
		db,
		s,
		enum.RoleAdmin,
	)
}

func addPartyRole(
	t *testing.T,
	db *gorm.DB,
	userID uuid.UUID,
	partyID uuid.UUID,
	role enum.PartyRole,
) appModels.PartyMember {

	member := appModels.PartyMember{
		ID: uuid.New(),

		UserID: userID,

		PartyID: partyID,
	}

	if err := db.Create(&member).Error; err != nil {
		t.Fatal(err)
	}

	partyRole := appModels.PartyMemberRole{
		ID: uuid.New(),

		PartyMemberID: member.ID,

		Role: role,
	}

	if err := db.Create(&partyRole).Error; err != nil {
		t.Fatal(err)
	}

	return member
}

func addSecondPartyStaff(
	t *testing.T,
	db *gorm.DB,
	partyID uuid.UUID,
) appModels.User {

	staff := fixtures.User()

	if err := db.Create(&staff).Error; err != nil {
		t.Fatal(err)
	}

	addPartyRole(
		t,
		db,
		staff.ID,
		partyID,
		enum.RoleStaff,
	)

	return staff
}

func createSecondParty(
	t *testing.T,
	db *gorm.DB,
	organizer appModels.User,
) appModels.Party {

	category := fixtures.Category()

	if err := db.Create(&category).Error; err != nil {
		t.Fatal(err)
	}

	party := fixtures.PartyWithOrganizer(
		organizer.ID,
	)

	party.CategoryID = category.ID

	if err := db.Create(&party).Error; err != nil {
		t.Fatal(err)
	}

	return party
}

func createAccessWindow(
	t *testing.T,
	db *gorm.DB,
	categoryID uuid.UUID,
	start time.Time,
	end time.Time,
) appModels.TicketAccessWindow {

	window := appModels.TicketAccessWindow{

		ID: uuid.New(),

		TicketCategoryID: categoryID,

		StartsAt: start,

		EndsAt: end,
	}

	if err := db.Create(&window).Error; err != nil {
		t.Fatal(err)
	}

	return window
}

func createNonVerificationCategory(
	t *testing.T,
	db *gorm.DB,
	partyID uuid.UUID,
) appModels.TicketCategory {

	category := appModels.TicketCategory{

		ID: uuid.New(),

		Name: "Standard",

		Price: 50,

		Capacity: 100,

		PartyID: partyID,

		RequiresVerification: false,
	}

	if err := db.Create(&category).Error; err != nil {
		t.Fatal(err)
	}

	return category
}

func createTicket(
	t *testing.T,
	db *gorm.DB,
	categoryID uuid.UUID,
	userID uuid.UUID,
	partyID uuid.UUID,
) appModels.Ticket {

	purchase := fixtures.Purchase(
		userID,
		partyID,
	)

	if err := db.Create(&purchase).Error; err != nil {
		t.Fatal(err)
	}

	ticket := fixtures.Ticket()

	ticket.TicketCategoryID = categoryID
	ticket.UserID = userID
	ticket.PurchaseID = purchase.ID

	if err := db.Create(&ticket).Error; err != nil {
		t.Fatal(err)
	}

	return ticket
}

func createPendingScan(
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

func assertScanStatus(
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

func verifyScan(
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

func assertVerified(
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

func assertRejected(
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

func changeTicketOwner(
	t *testing.T,
	db *gorm.DB,
	ticket *appModels.Ticket,
	userID uuid.UUID,
) {

	if err := db.
		Model(ticket).
		Update(
			"user_id",
			userID,
		).
		Error; err != nil {

		t.Fatal(err)
	}
}
