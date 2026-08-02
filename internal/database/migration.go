package database

import (
	"gorm.io/gorm"

	"github.com/reinp/event-platform/backend/internal/models"
)

func Migrate(
	db *gorm.DB,
) error {
	if err := migrateLegacyUUIDColumns(db); err != nil {
		return err
	}

	if err := db.AutoMigrate(

		// auth
		&models.User{},
		&models.UserRole{},
		&models.RefreshToken{},
		&models.PasswordResetToken{},
		&models.EmailVerification{},

		// media
		&models.Media{},

		// artist system
		&models.Artist{},

		// party system
		&models.PartyCategory{},
		&models.Party{},
		&models.PartyMedia{},

		// party permissions
		&models.PartyMember{},
		&models.PartyMemberRole{},

		// organizational staff
		&models.StaffMember{},

		// lineup system
		&models.PartyStage{},
		&models.ArtistSlot{},

		// ticketing
		&models.TicketCategory{},
		&models.TicketAccessWindow{},
		&models.Ticket{},
		&models.TicketScan{},

		// purchases
		&models.Purchase{},
		&models.PurchaseItem{},

		// payments
		&models.PaymentEvent{},
	); err != nil {
		return err
	}

	return migrateLegacyUserRoles(db)
}

func migrateLegacyUUIDColumns(db *gorm.DB) error {
	return db.Exec(`
		DO $$
		BEGIN
			IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_schema=current_schema() AND table_name='artists' AND column_name='id' AND data_type='text') THEN
				ALTER TABLE artists ALTER COLUMN id DROP DEFAULT;
				ALTER TABLE artists ALTER COLUMN id TYPE uuid USING id::uuid;
				ALTER TABLE artists ALTER COLUMN id SET DEFAULT gen_random_uuid();
			END IF;
			IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_schema=current_schema() AND table_name='party_stages' AND column_name='id' AND data_type='text') THEN
				ALTER TABLE party_stages ALTER COLUMN id DROP DEFAULT;
				ALTER TABLE party_stages ALTER COLUMN id TYPE uuid USING id::uuid;
				ALTER TABLE party_stages ALTER COLUMN id SET DEFAULT gen_random_uuid();
			END IF;
			IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_schema=current_schema() AND table_name='artist_slots' AND column_name='id' AND data_type='text') THEN
				ALTER TABLE artist_slots ALTER COLUMN id DROP DEFAULT;
				ALTER TABLE artist_slots ALTER COLUMN id TYPE uuid USING id::uuid;
				ALTER TABLE artist_slots ALTER COLUMN id SET DEFAULT gen_random_uuid();
			END IF;
			IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_schema=current_schema() AND table_name='staff_members' AND column_name='id' AND data_type='text') THEN
				ALTER TABLE staff_members ALTER COLUMN id DROP DEFAULT;
				ALTER TABLE staff_members ALTER COLUMN id TYPE uuid USING id::uuid;
				ALTER TABLE staff_members ALTER COLUMN id SET DEFAULT gen_random_uuid();
			END IF;
			IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_schema=current_schema() AND table_name='ticket_categories' AND column_name='refund_policy_id' AND data_type='text') THEN
				ALTER TABLE ticket_categories ALTER COLUMN refund_policy_id TYPE uuid USING refund_policy_id::uuid;
			END IF;
		END $$;
	`).Error
}

func migrateLegacyUserRoles(db *gorm.DB) error {
	return db.Exec(`
		DO $$
		BEGIN
			IF EXISTS (
				SELECT 1
				FROM information_schema.columns
				WHERE table_schema = current_schema()
				  AND table_name = 'users'
				  AND column_name = 'role'
			) THEN
				INSERT INTO user_roles (id, user_id, role, created_at)
				SELECT gen_random_uuid(), id, role, NOW()
				FROM users
				WHERE role IS NOT NULL
				ON CONFLICT (user_id, role) DO NOTHING;
			END IF;
		END $$;
	`).Error
}
