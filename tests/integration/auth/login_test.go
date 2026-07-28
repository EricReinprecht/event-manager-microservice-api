package auth_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/reinp/event-platform/backend/tests/helpers"
	"github.com/reinp/event-platform/backend/tests/scenarios"
)

func TestLoginSuccess(
	t *testing.T,
) {

	ctx := context.Background()

	db := setupAuthTest(
		t,
	)

	authService := helpers.NewAuthService(
		db,
	)

	scenario := scenarios.NewAuthScenario(
		db,
		authService,
	)

	user := scenario.CreateVerifiedUser(
		ctx,
	)

	response := scenario.Login(
		ctx,
		user,
	)

	require.NotNil(
		t,
		response,
	)

	assert.NotEmpty(
		t,
		response.AccessToken,
	)

	assert.NotEmpty(
		t,
		response.RefreshToken,
	)

	var count int64

	err := db.
		Table("refresh_tokens").
		Where(
			"user_id = ?",
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

	assert.Equal(
		t,
		int64(1),
		count,
	)
}
