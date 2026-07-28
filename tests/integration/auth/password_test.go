package auth_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/reinp/event-platform/backend/internal/auth"
	"github.com/reinp/event-platform/backend/internal/database"
	"github.com/reinp/event-platform/backend/internal/models"
	"github.com/reinp/event-platform/backend/internal/service"

	"github.com/reinp/event-platform/backend/tests/helpers"
	"github.com/reinp/event-platform/backend/tests/scenarios"
)

func TestChangePasswordSuccess(t *testing.T) {

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

func TestChangePasswordWrongOldPassword(t *testing.T) {

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

func TestChangePasswordValidatorRejectsWeakPassword(t *testing.T) {

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

func TestChangePasswordCreatesValidHash(t *testing.T) {

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

func TestChangePasswordRevokesAllSessions(t *testing.T) {

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

	// create multiple sessions
	login1 := scenario.Login(
		ctx,
		user,
	)

	login2 := scenario.Login(
		ctx,
		user,
	)

	// change password
	err := userService.ChangePassword(
		ctx,
		user.ID,
		"C4ctus!River#829Lamp",
		"V7q!Nexa42#Orbit91",
	)

	require.NoError(
		t,
		err,
	)

	// all refresh tokens should be revoked

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

	// old refresh tokens should not work anymore

	_, err = authService.Refresh(
		ctx,
		login1.RefreshToken,
	)

	require.Error(
		t,
		err,
	)

	_, err = authService.Refresh(
		ctx,
		login2.RefreshToken,
	)

	require.Error(
		t,
		err,
	)
}

func TestPasswordResetSuccess(t *testing.T) {

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

	// create active session before password reset

	_ = scenario.Login(
		ctx,
		user,
	)

	executor := database.NewGormExecutor(db)

	rawToken, err := helpers.CreatePasswordResetToken(
		executor,
		user.ID,
	)

	require.NoError(
		t,
		err,
	)

	newPassword := "C4ctus!River#829Lamp"

	err = authService.ResetPassword(
		ctx,
		rawToken,
		newPassword,
	)

	require.NoError(
		t,
		err,
	)

	// old password fails

	_, err = authService.Login(
		ctx,
		user.Email,
		"Password123!",
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
		newPassword,
		"test",
		"127.0.0.1",
	)

	require.NoError(
		t,
		err,
	)

	// existing sessions were revoked

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

func TestPasswordResetInvalidToken(t *testing.T) {

	ctx := context.Background()

	db := setupAuthTest(t)

	authService := helpers.NewAuthService(
		db,
	)

	err := authService.ResetPassword(
		ctx,
		"invalid-reset-token",
		"C4ctus!River#829Lamp",
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

func TestPasswordResetExpiredToken(t *testing.T) {

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

	executor := database.NewGormExecutor(db)

	rawToken, err := helpers.CreatePasswordResetToken(
		executor,
		user.ID,
	)

	require.NoError(
		t,
		err,
	)

	// manually expire token in database
	err = db.
		Model(
			&models.PasswordResetToken{},
		).
		Where(
			"user_id = ?",
			user.ID,
		).
		Update(
			"expires_at",
			time.Now().Add(-1*time.Hour),
		).
		Error

	require.NoError(
		t,
		err,
	)

	err = authService.ResetPassword(
		ctx,
		rawToken,
		"Violet!Forest#839Tree",
	)

	require.Error(
		t,
		err,
	)

	assert.Contains(
		t,
		err.Error(),
		"expired",
	)
}

func TestPasswordResetCannotReuseToken(t *testing.T) {

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

	executor := database.NewGormExecutor(db)

	rawToken, err := helpers.CreatePasswordResetToken(
		executor,
		user.ID,
	)

	require.NoError(
		t,
		err,
	)

	// first reset should succeed

	err = authService.ResetPassword(
		ctx,
		rawToken,
		"Violet!Forest#839Tree",
	)

	require.NoError(
		t,
		err,
	)

	// second reset with same token must fail

	err = authService.ResetPassword(
		ctx,
		rawToken,
		"Another!Password#839Tree",
	)

	require.Error(
		t,
		err,
	)

	assert.Contains(
		t,
		err.Error(),
		"used",
	)
}

func TestPasswordResetInvalidatesOtherResetTokens(t *testing.T) {

	ctx := context.Background()

	db := setupAuthTest(t)

	authService := helpers.NewAuthService(
		db,
	)

	req := service.RegisterRequest{
		Email: "reset-invalidate-" + uuid.NewString() + "@example.com",

		Username: "resetinvalidate" + uuid.NewString()[:8],

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

	// create first reset token
	rawToken1 := auth.GenerateToken()

	reset1 := &models.PasswordResetToken{
		UserID: user.ID,

		TokenHash: auth.HashToken(
			rawToken1,
		),

		ExpiresAt: time.Now().Add(
			15 * time.Minute,
		),
	}

	err = db.Create(
		reset1,
	).Error

	require.NoError(
		t,
		err,
	)

	// create second reset token (simulates another reset request)
	err = authService.ForgotPassword(
		ctx,
		req.Email,
	)

	require.NoError(
		t,
		err,
	)

	// old token should now be invalidated

	err = authService.ResetPassword(
		ctx,
		rawToken1,
		"NewStrongPassword!12345",
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

func TestPasswordResetRevokesAllSessions(t *testing.T) {

	ctx := context.Background()

	db := setupAuthTest(t)

	authService := helpers.NewAuthService(
		db,
	)

	req := service.RegisterRequest{
		Email: "reset-session-" + uuid.NewString() + "@example.com",

		Username: "resetsession" + uuid.NewString()[:8],

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

	// verify user

	now := time.Now()

	user.VerifiedAt = &now

	err = db.Save(
		user,
	).Error

	require.NoError(
		t,
		err,
	)

	// create active session

	tokens, err := authService.Login(
		ctx,
		req.Email,
		req.Password,
		"test-device",
		"127.0.0.1",
	)

	require.NoError(
		t,
		err,
	)

	require.NotEmpty(
		t,
		tokens.RefreshToken,
	)

	// create password reset token

	rawResetToken := auth.GenerateToken()

	reset := &models.PasswordResetToken{

		UserID: user.ID,

		TokenHash: auth.HashToken(
			rawResetToken,
		),

		ExpiresAt: time.Now().Add(
			15 * time.Minute,
		),
	}

	err = db.Create(
		reset,
	).Error

	require.NoError(
		t,
		err,
	)

	// reset password

	err = authService.ResetPassword(
		ctx,
		rawResetToken,
		"C4ctus!River#829LamsdpX",
	)

	require.NoError(
		t,
		err,
	)

	// old refresh token must no longer work

	_, err = authService.Refresh(
		ctx,
		tokens.RefreshToken,
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

func TestPasswordResetCannotResetDeletedUser(t *testing.T) {

	ctx := context.Background()

	db := setupAuthTest(t)

	emailService := &helpers.CapturingEmailService{}

	authService := helpers.NewAuthServiceWithCapturingEmail(
		db,
		emailService,
	)

	userService := helpers.NewUserService(db)

	// create user
	user, err := authService.Register(
		ctx,
		service.RegisterRequest{
			Email:    "deletedreset@example.com",
			Username: "deletedreset",
			Password: "C4ctus!River#829Lamp",
		},
	)

	require.NoError(
		t,
		err,
	)

	// verify email
	verificationToken := emailService.VerificationToken

	require.NotEmpty(
		t,
		verificationToken,
	)

	_, err = authService.VerifyEmail(
		ctx,
		verificationToken,
	)

	require.NoError(
		t,
		err,
	)

	// create reset token
	err = authService.ForgotPassword(
		ctx,
		user.Email,
	)

	require.NoError(
		t,
		err,
	)

	resetToken := emailService.ResetToken

	require.NotEmpty(
		t,
		resetToken,
	)

	// delete user
	err = userService.Delete(
		ctx,
		user.ID,
	)

	require.NoError(
		t,
		err,
	)

	// try reset after deletion
	err = authService.ResetPassword(
		ctx,
		resetToken,
		"NewStrongPassword!123",
	)

	require.Error(
		t,
		err,
	)
}
