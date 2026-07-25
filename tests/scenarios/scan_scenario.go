package scenarios

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/reinp/event-platform/backend/internal/models"
	appModels "github.com/reinp/event-platform/backend/internal/models"
	"github.com/reinp/event-platform/backend/internal/models/enum"
	"github.com/reinp/event-platform/backend/internal/service"
	"github.com/reinp/event-platform/backend/tests/fixtures"
	"github.com/reinp/event-platform/backend/tests/helpers"
	"gorm.io/gorm"
)

type ScanScenario struct {
	DB *gorm.DB

	TicketService *service.TicketService

	Staff models.User

	StaffTwo models.User

	Customer models.User

	Organizer models.User

	Party models.Party

	Category models.Category

	OrganizerMember models.PartyMember

	Member models.PartyMember

	MemberTwo models.PartyMember

	TicketCategory models.TicketCategory

	Window models.TicketAccessWindow

	Ticket models.Ticket
}

func CreateScanScenario(
	t *testing.T,
	db *gorm.DB,
	clock *helpers.FakeClock,
	options ...bool,
) ScanScenario {

	t.Helper()

	requiresVerification := false

	if len(options) > 0 {
		requiresVerification = options[0]
	}

	ticketService := helpers.NewTicketService(
		db,
		clock,
	)

	// USERS

	staff := fixtures.User()

	if err := db.Create(&staff).Error; err != nil {
		t.Fatal(err)
	}

	staffTwo := fixtures.User()

	if err := db.Create(&staffTwo).Error; err != nil {
		t.Fatal(err)
	}

	customer := fixtures.User()

	if err := db.Create(&customer).Error; err != nil {
		t.Fatal(err)
	}

	// CATEGORY

	category := fixtures.Category()

	if err := db.Create(&category).Error; err != nil {
		t.Fatal(err)
	}

	// ORGANIZER

	organizer := fixtures.User()

	if err := db.Create(&organizer).Error; err != nil {
		t.Fatal(err)
	}

	// PARTY

	party := fixtures.PartyWithOrganizer(
		organizer.ID,
	)

	party.CategoryID = category.ID

	if err := db.Create(&party).Error; err != nil {
		t.Fatal(err)
	}

	organizerMember := AddPartyRole(
		t,
		db,
		organizer.ID,
		party.ID,
		enum.RoleOrganizer,
	)

	// IMPORTANT:
	// PartyWithOrganizer already created organizer membership.
	// Do NOT add another membership for staff here.

	member := AddPartyRole(
		t,
		db,
		staff.ID,
		party.ID,
		enum.RoleStaff,
	)

	memberTwo := AddPartyRole(
		t,
		db,
		staffTwo.ID,
		party.ID,
		enum.RoleStaff,
	)
	// TICKET CATEGORY

	ticketCategory := models.TicketCategory{

		ID: uuid.New(),

		Name: "VIP",

		PartyID: party.ID,

		RequiresVerification: requiresVerification,
	}

	if err := db.Create(&ticketCategory).Error; err != nil {
		t.Fatal(err)
	}

	// ACCESS WINDOW

	window := models.TicketAccessWindow{

		ID: uuid.New(),

		TicketCategoryID: ticketCategory.ID,

		StartsAt: clock.Now().Add(-time.Hour),

		EndsAt: clock.Now().Add(time.Hour),
	}

	if err := db.Create(&window).Error; err != nil {
		t.Fatal(err)
	}

	// PURCHASE

	purchase := CreatePurchase(
		t,
		db,
		customer.ID,
		party.ID,
	)

	// TICKET

	ticket := fixtures.Ticket()

	ticket.TicketCategoryID = ticketCategory.ID

	ticket.UserID = customer.ID

	ticket.PurchaseID = purchase.ID

	if err := db.Create(&ticket).Error; err != nil {
		t.Fatal(err)
	}

	return ScanScenario{

		DB: db,

		TicketService: ticketService,

		Organizer: organizer,

		OrganizerMember: organizerMember,

		Staff: staff,

		StaffTwo: staffTwo,

		Customer: customer,

		Party: party,

		Category: category,

		Member: member,

		MemberTwo: memberTwo,

		TicketCategory: ticketCategory,

		Window: window,

		Ticket: ticket,
	}
}

func CreateAdditionalTicket(
	t *testing.T,
	db *gorm.DB,
	categoryID uuid.UUID,
	userID uuid.UUID,
	partyID uuid.UUID,
) models.Ticket {

	purchase := CreatePurchase(
		t,
		db,
		userID,
		partyID,
	)

	ticket := fixtures.Ticket()

	ticket.TicketCategoryID = categoryID

	ticket.UserID = userID

	ticket.PurchaseID = purchase.ID

	if err := db.Create(&ticket).Error; err != nil {
		t.Fatal(err)
	}

	return ticket
}

func CreatePendingTicketScan(
	t *testing.T,
	db *gorm.DB,
	ticket models.Ticket,
	window models.TicketAccessWindow,
	scanner models.User,
	at time.Time,
) models.TicketScan {

	t.Helper()

	scan := models.TicketScan{

		ID: uuid.New(),

		TicketID: ticket.ID,

		TicketAccessWindowID: window.ID,

		ScannedByID: scanner.ID,

		ScannedAt: at,

		Status: enum.TicketScanPending,
	}

	if err := db.Create(&scan).Error; err != nil {
		t.Fatal(err)
	}

	return scan
}

func CancelTicket(
	t *testing.T,
	db *gorm.DB,
	ticket *models.Ticket,
) {

	t.Helper()

	ticket.Status = enum.TicketStatusCancelled

	if err := db.Save(ticket).Error; err != nil {
		t.Fatal(err)
	}
}

func AddOrganizerRole(
	t *testing.T,
	db *gorm.DB,
	userID uuid.UUID,
	partyID uuid.UUID,
) appModels.PartyMember {

	return AddPartyRole(
		t,
		db,
		userID,
		partyID,
		enum.RoleOrganizer,
	)
}

func RemovePartyMembership(
	t *testing.T,
	db *gorm.DB,
	member *models.PartyMember,
) {
	RemovePartyMember(
		t,
		db,
		member,
	)
}

func CreateStaffParty(
	t *testing.T,
	db *gorm.DB,
	userID uuid.UUID,
	categoryID uuid.UUID,
) models.Party {

	t.Helper()

	organizer := fixtures.User()

	if err := db.Create(&organizer).Error; err != nil {
		t.Fatal(err)
	}

	party := fixtures.PartyWithOrganizer(
		organizer.ID,
	)

	party.CategoryID = categoryID

	if err := db.Create(&party).Error; err != nil {
		t.Fatal(err)
	}

	AddPartyRole(
		t,
		db,
		userID,
		party.ID,
		enum.RoleStaff,
	)

	return party
}

func CreateStaffInParty(
	t *testing.T,
	db *gorm.DB,
	user models.User,
	category models.Category,
) models.Party {

	// ORGANIZER

	organizer := fixtures.User()

	if err := db.Create(&organizer).Error; err != nil {
		t.Fatal(err)
	}

	// PARTY

	party := fixtures.PartyWithOrganizer(
		organizer.ID,
	)

	party.CategoryID = category.ID

	if err := db.Create(&party).Error; err != nil {
		t.Fatal(err)
	}

	AddPartyRole(
		t,
		db,
		user.ID,
		party.ID,
		enum.RoleStaff,
	)

	return party
}
