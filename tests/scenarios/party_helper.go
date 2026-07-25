package scenarios

import (
	"testing"

	"github.com/google/uuid"

	appModels "github.com/reinp/event-platform/backend/internal/models"
	"github.com/reinp/event-platform/backend/tests/fixtures"

	"github.com/reinp/event-platform/backend/internal/models/enum"
	"gorm.io/gorm"
)

func AddPartyRole(
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

func MakeOrganizer(
	t *testing.T,
	db *gorm.DB,
	s *VerificationScenario,
) {
	ChangeStaffRole(
		t,
		db,
		s,
		enum.RoleOrganizer,
	)
}

func MakeAdmin(
	t *testing.T,
	db *gorm.DB,
	s *VerificationScenario,
) {
	ChangeStaffRole(
		t,
		db,
		s,
		enum.RoleAdmin,
	)
}

func AddSecondPartyStaff(
	t *testing.T,
	db *gorm.DB,
	partyID uuid.UUID,
) appModels.User {

	staff := fixtures.User()

	if err := db.Create(&staff).Error; err != nil {
		t.Fatal(err)
	}

	AddPartyRole(
		t,
		db,
		staff.ID,
		partyID,
		enum.RoleStaff,
	)

	return staff
}

func CreateSecondParty(
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

func ChangeStaffRole(
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
