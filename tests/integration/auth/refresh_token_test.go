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
	"github.com/reinp/event-platform/backend/internal/service"

	"github.com/reinp/event-platform/backend/tests/helpers"
	"github.com/reinp/event-platform/backend/tests/scenarios"
)

func TestRefreshTokenSuccess(t *testing.T) {

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

func TestRefreshTokenReplayDetection(t *testing.T) {

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

func TestRefreshTokenExpired(t *testing.T) {

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

func TestRefreshTokenAfterLogoutAll(t *testing.T) {

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

	require.NotEmpty(
		t,
		login.RefreshToken,
	)

	// revoke all sessions

	err := authService.LogoutAll(
		ctx,
		user.ID,
	)

	require.NoError(
		t,
		err,
	)

	// refresh token should no longer work

	_, err = authService.Refresh(
		ctx,
		login.RefreshToken,
	)

	require.Error(
		t,
		err,
	)

	assert.Contains(
		t,
		err.Error(),
		"replay",
	)
}

func TestRefreshTokenAfterPasswordChange(t *testing.T) {

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

	login := scenario.Login(
		ctx,
		user,
	)

	require.NotEmpty(
		t,
		login.RefreshToken,
	)

	err := userService.ChangePassword(
		ctx,
		user.ID,
		"C4ctus!River#829Lamp",
		"Violet!Forest#839Tree",
	)

	require.NoError(
		t,
		err,
	)

	// old refresh token should no longer work

	_, err = authService.Refresh(
		ctx,
		login.RefreshToken,
	)

	require.Error(
		t,
		err,
	)

	assert.Contains(
		t,
		err.Error(),
		"replay",
	)
}

func TestDeletedUserCannotRefreshToken(t *testing.T) {

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

	require.NotEmpty(
		t,
		login.RefreshToken,
	)

	// delete user

	err := db.
		Delete(
			&user,
		).
		Error

	require.NoError(
		t,
		err,
	)

	// refresh should fail

	_, err = authService.Refresh(
		ctx,
		login.RefreshToken,
	)

	require.Error(
		t,
		err,
	)

	assert.Contains(
		t,
		err.Error(),
		"invalid",
	)
}

func TestRefreshTokenCannotBeUsedTwice(t *testing.T) {

	ctx := context.Background()

	db := setupAuthTest(t)

	authService := helpers.NewAuthService(
		db,
	)

	// create verified user
	req := service.RegisterRequest{
		Email: "refresh-double-" + uuid.NewString() + "@example.com",

		Username: "refreshtwice" + uuid.NewString()[:8],

		Password: "C4ctus!River#829Lamp",
	}

	user, err := authService.Register(
		ctx,
		req,
	)

	require.NoError(
		t,
		err,
	)

	now := time.Now()

	user.VerifiedAt = &now

	err = db.Save(user).Error

	require.NoError(
		t,
		err,
	)

	// login -> get refresh token
	tokens, err := authService.Login(
		ctx,
		user.Email,
		req.Password,
		"test-device",
		"127.0.0.1",
	)

	require.NoError(
		t,
		err,
	)

	oldRefreshToken := tokens.RefreshToken

	require.NotEmpty(
		t,
		oldRefreshToken,
	)

	// first refresh succeeds
	newTokens, err := authService.Refresh(
		ctx,
		oldRefreshToken,
	)

	require.NoError(
		t,
		err,
	)

	require.NotEmpty(
		t,
		newTokens.RefreshToken,
	)

	// old token cannot be reused
	_, err = authService.Refresh(
		ctx,
		oldRefreshToken,
	)

	require.Error(
		t,
		err,
	)

	assert.Contains(
		t,
		err.Error(),
		"replay",
	)
}

func TestRefreshTokenRotation(t *testing.T) {

	ctx := context.Background()

	db := setupAuthTest(t)

	authService := helpers.NewAuthService(
		db,
	)

	user, err := authService.Register(
		ctx,
		service.RegisterRequest{
			Email: "rotation-" + uuid.NewString() + "@example.com",

			Username: "rotationuser" + uuid.NewString()[:8],

			Password: "C4ctus!River#829Lamp",
		},
	)

	require.NoError(
		t,
		err,
	)

	now := time.Now()

	user.VerifiedAt = &now

	require.NoError(
		t,
		db.Save(user).Error,
	)

	// login

	firstTokens, err := authService.Login(
		ctx,
		user.Email,
		"C4ctus!River#829Lamp",
		"browser",
		"127.0.0.1",
	)

	require.NoError(
		t,
		err,
	)

	oldRefreshToken := firstTokens.RefreshToken

	// refresh

	secondTokens, err := authService.Refresh(
		ctx,
		oldRefreshToken,
	)

	require.NoError(
		t,
		err,
	)

	require.NotEmpty(
		t,
		secondTokens.RefreshToken,
	)

	assert.NotEqual(
		t,
		oldRefreshToken,
		secondTokens.RefreshToken,
	)

	// old token cannot be reused

	_, err = authService.Refresh(
		ctx,
		oldRefreshToken,
	)

	require.Error(
		t,
		err,
	)

	assert.Contains(
		t,
		err.Error(),
		"replay",
	)
}

func TestRefreshTokenFamilyIsolation(t *testing.T) {

	ctx := context.Background()

	db := setupAuthTest(t)

	authService := helpers.NewAuthService(
		db,
	)

	user, err := authService.Register(
		ctx,
		service.RegisterRequest{
			Email: "family-isolation-" + uuid.NewString() + "@example.com",

			Username: "familyuser" + uuid.NewString()[:8],

			Password: "C4ctus!River#829Lamp",
		},
	)

	require.NoError(
		t,
		err,
	)

	now := time.Now()

	user.VerifiedAt = &now

	require.NoError(
		t,
		db.Save(user).Error,
	)

	// first login/session

	sessionA, err := authService.Login(
		ctx,
		user.Email,
		"C4ctus!River#829Lamp",
		"browser-a",
		"127.0.0.1",
	)

	require.NoError(
		t,
		err,
	)

	// second login/session

	sessionB, err := authService.Login(
		ctx,
		user.Email,
		"C4ctus!River#829Lamp",
		"browser-b",
		"127.0.0.2",
	)

	require.NoError(
		t,
		err,
	)

	require.NotEqual(
		t,
		sessionA.RefreshToken,
		sessionB.RefreshToken,
	)

	// find family A

	var tokenA models.RefreshToken

	err = db.
		Where(
			"token_hash = ?",
			auth.HashToken(sessionA.RefreshToken),
		).
		First(
			&tokenA,
		).
		Error

	require.NoError(
		t,
		err,
	)

	// revoke only family A

	err = authService.RevokeSession(
		ctx,
		user.ID,
		tokenA.FamilyID,
	)

	require.NoError(
		t,
		err,
	)

	// family A should fail

	_, err = authService.Refresh(
		ctx,
		sessionA.RefreshToken,
	)

	require.Error(
		t,
		err,
	)

	// family B should still work

	newTokens, err := authService.Refresh(
		ctx,
		sessionB.RefreshToken,
	)

	require.NoError(
		t,
		err,
	)

	require.NotNil(
		t,
		newTokens,
	)

	assert.NotEmpty(
		t,
		newTokens.AccessToken,
	)
}

func TestConcurrentRefreshRace(t *testing.T) {

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

	tokens := scenario.Login(
		ctx,
		user,
	)

	refreshToken := tokens.RefreshToken

	require.NotEmpty(
		t,
		refreshToken,
	)

	type result struct {
		err error
	}

	results := make(chan result, 2)

	// run two refreshes at the same time

	go func() {

		_, err := authService.Refresh(
			context.Background(),
			refreshToken,
		)

		results <- result{
			err: err,
		}
	}()

	go func() {

		_, err := authService.Refresh(
			context.Background(),
			refreshToken,
		)

		results <- result{
			err: err,
		}
	}()

	first := <-results
	second := <-results

	successCount := 0
	errorCount := 0

	for _, r := range []result{first, second} {

		if r.err == nil {
			successCount++
		} else {
			errorCount++
		}
	}

	// exactly one refresh should succeed

	assert.Equal(
		t,
		1,
		successCount,
	)

	assert.Equal(
		t,
		1,
		errorCount,
	)
}
