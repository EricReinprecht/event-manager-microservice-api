package helpers

import (
	"sync"

	"gorm.io/gorm"
)

var DatabaseCleanupMutex sync.Mutex

func CleanDatabase(db *gorm.DB) error {

	DatabaseCleanupMutex.Lock()
	defer DatabaseCleanupMutex.Unlock()

	return db.Exec(`
		TRUNCATE TABLE
			ticket_scans,
			payment_events,
			tickets,
			purchase_items,
			purchases,
			ticket_access_windows,
			ticket_categories,
			party_member_roles,
			party_members,
			parties,
			categories,
			users
		RESTART IDENTITY CASCADE
	`).Error
}
