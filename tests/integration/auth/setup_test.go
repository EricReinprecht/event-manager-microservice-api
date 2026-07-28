package auth_test

import (
	"testing"

	"gorm.io/gorm"

	"github.com/reinp/event-platform/backend/tests/helpers"
)

func setupAuthTest(
	t *testing.T,
) *gorm.DB {

	t.Helper()

	db, err := helpers.TestDatabase()

	if err != nil {
		t.Fatal(err)
	}

	if err := helpers.CleanDatabase(db); err != nil {
		t.Fatal(err)
	}

	return db
}
