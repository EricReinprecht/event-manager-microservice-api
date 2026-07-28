package auth_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/reinp/event-platform/backend/internal/auth"
	"github.com/reinp/event-platform/backend/internal/models"

	"github.com/reinp/event-platform/backend/tests/helpers"
	"github.com/reinp/event-platform/backend/tests/scenarios"
)

func TestRefreshTokenSuccess(
	t *testing.T,
) {

	ctx := context.Background()

	db := setupAuthTest(t)

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

	login := scenario.Login(
		ctx,
		user,
	)

	result, err := authService.Refresh(
		ctx,
		login.RefreshToken,
	)

	require.NoError(
		t,
		err,
	)

	assert.NotEmpty(
		t,
		result.AccessToken,
	)

	assert.NotEmpty(
		t,
		result.RefreshToken,
	)

	assert.NotEqual(
		t,
		login.RefreshToken,
		result.RefreshToken,
	)

	var revokedToken models.RefreshToken

	err = db.
		Where(
			"user_id = ? AND revoked_at IS NOT NULL",
			user.ID,
		).
		First(
			&revokedToken,
		).
		Error

	require.NoError(
		t,
		err,
	)

	assert.NotNil(
		t,
		revokedToken.RevokedAt,
	)
}

func TestRefreshTokenReplayDetection(
	t *testing.T,
) {

	ctx := context.Background()

	db := setupAuthTest(t)

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

	login := scenario.Login(
		ctx,
		user,
	)

	_, err := authService.Refresh(
		ctx,
		login.RefreshToken,
	)

	require.NoError(
		t,
		err,
	)

	_, err = authService.Refresh(
		ctx,
		login.RefreshToken,
	)

	require.Error(
		t,
		err,
	)

	assert.Equal(
		t,
		"refresh token replay detected",
		err.Error(),
	)

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
		int64(2),
	)
}

func TestRefreshTokenExpired(
	t *testing.T,
) {

	ctx := context.Background()

	db := setupAuthTest(t)

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

	rawToken := auth.GenerateToken()

	token := &models.RefreshToken{

		UserID: user.ID,

		TokenHash: auth.HashToken(
			rawToken,
		),

		FamilyID: uuid.New(),

		UserAgent: "test",

		IPAddress: "127.0.0.1",

		DeviceName: "test",

		ExpiresAt: time.Now().Add(
			-1 * time.Hour,
		),
	}

	err := db.Create(
		token,
	).Error

	require.NoError(
		t,
		err,
	)

	_, err = authService.Refresh(
		ctx,
		rawToken,
	)

	require.Error(
		t,
		err,
	)

	assert.Equal(
		t,
		"refresh token expired",
		err.Error(),
	)
}
