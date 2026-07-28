package auth_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/reinp/event-platform/backend/internal/models"

	"github.com/reinp/event-platform/backend/tests/helpers"
	"github.com/reinp/event-platform/backend/tests/scenarios"
)

func TestChangePasswordSuccess(
	t *testing.T,
) {

	ctx := context.Background()

	db := setupAuthTest(t)

	authService := helpers.NewAuthService(
		db,
	)

	userService := helpers.NewUserService(db)

	scenario := scenarios.NewAuthScenario(
		db,
		authService,
	)

	user := scenario.CreateVerifiedUser(
		ctx,
	)

	scenario.Login(
		ctx,
		user,
	)

	err := userService.ChangePassword(
		ctx,
		user.ID,
		"C4ctus!River#829Lamp",
		"NewC4ctus!River#829Lamp",
	)

	require.NoError(
		t,
		err,
	)

	// old password no longer works

	_, err = authService.Login(
		ctx,
		user.Email,
		"C4ctus!River#829Lamp",
		"test",
		"127.0.0.1",
	)

	require.Error(
		t,
		err,
	)

	// new password works

	_, err = authService.Login(
		ctx,
		user.Email,
		"NewC4ctus!River#829Lamp",
		"test",
		"127.0.0.1",
	)

	require.NoError(
		t,
		err,
	)

	// old sessions should be revoked

	var count int64

	err = db.
		Table(
			"refresh_tokens",
		).
		Where(
			"user_id = ? AND revoked_at IS NOT NULL",
			user.ID,
		).
		Count(
			&count,
		).
		Error

	require.NoError(
		t,
		err,
	)

	assert.GreaterOrEqual(
		t,
		count,
		int64(1),
	)
}

func TestChangePasswordWrongOldPassword(
	t *testing.T,
) {

	ctx := context.Background()

	db := setupAuthTest(t)

	authService := helpers.NewAuthService(
		db,
	)

	userService := helpers.NewUserService(
		db,
	)

	scenario := scenarios.NewAuthScenario(
		db,
		authService,
	)

	user := scenario.CreateVerifiedUser(
		ctx,
	)

	err := userService.ChangePassword(
		ctx,
		user.ID,
		"WrongC4ctus!River#829Lamp",
		"NewC4ctus!River#829Lamp",
	)

	require.Error(
		t,
		err,
	)

	assert.Equal(
		t,
		"invalid old password",
		err.Error(),
	)
}

func TestChangePasswordValidatorRejectsWeakPassword(
	t *testing.T,
) {

	ctx := context.Background()

	db := setupAuthTest(t)

	authService := helpers.NewAuthService(
		db,
	)

	userService := helpers.NewUserService(
		db,
	)

	scenario := scenarios.NewAuthScenario(
		db,
		authService,
	)

	user := scenario.CreateVerifiedUser(
		ctx,
	)

	err := userService.ChangePassword(
		ctx,
		user.ID,
		"C4ctus!River#829Lamp",
		"123",
	)

	require.Error(
		t,
		err,
	)
}

func TestChangePasswordCreatesValidHash(
	t *testing.T,
) {

	ctx := context.Background()

	db := setupAuthTest(t)

	authService := helpers.NewAuthService(
		db,
	)

	userService := helpers.NewUserService(
		db,
	)

	scenario := scenarios.NewAuthScenario(
		db,
		authService,
	)

	user := scenario.CreateVerifiedUser(
		ctx,
	)

	oldHash := user.PasswordHash

	err := userService.ChangePassword(
		ctx,
		user.ID,
		"C4ctus!River#829Lamp",
		"AnotherC4ctus!River#829Lamp",
	)

	require.NoError(
		t,
		err,
	)

	var updated models.User

	err = db.
		Where(
			"id = ?",
			user.ID,
		).
		First(
			&updated,
		).
		Error

	require.NoError(
		t,
		err,
	)

	assert.NotEqual(
		t,
		oldHash,
		updated.PasswordHash,
	)
}
