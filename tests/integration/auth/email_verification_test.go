package auth_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/reinp/event-platform/backend/internal/models"
	"github.com/reinp/event-platform/backend/internal/service"

	"github.com/reinp/event-platform/backend/tests/helpers"
)

func TestEmailVerificationSuccess(t *testing.T) {

	ctx := context.Background()

	db := setupAuthTest(t)

	authService := helpers.NewAuthService(
		db,
	)

	req := service.RegisterRequest{
		Email: "verify-" + uuid.NewString() + "@example.com",

		Username: "verifysuccess" + uuid.NewString()[:8],

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

	rawToken, verification :=
		helpers.CreateEmailVerification(
			user.ID,
		)

	err = db.Create(
		verification,
	).Error

	require.NoError(
		t,
		err,
	)

	_, err = authService.VerifyEmail(
		ctx,
		rawToken,
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

	assert.NotNil(
		t,
		updated.VerifiedAt,
	)
}

func TestEmailVerificationInvalidToken(t *testing.T) {

	ctx := context.Background()

	db := setupAuthTest(t)

	authService := helpers.NewAuthService(
		db,
	)

	_, err := authService.VerifyEmail(
		ctx,
		"invalid-token",
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

func TestEmailVerificationAlreadyUsed(t *testing.T) {

	ctx := context.Background()

	db := setupAuthTest(t)

	authService := helpers.NewAuthService(
		db,
	)

	req := service.RegisterRequest{
		Email: "verify-used-" + uuid.NewString() + "@example.com",

		Username: "verifyused" + uuid.NewString()[:8],

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

	rawToken, verification :=
		helpers.CreateEmailVerification(
			user.ID,
		)

	err = db.Create(
		verification,
	).Error

	require.NoError(
		t,
		err,
	)

	// first verification succeeds

	_, err = authService.VerifyEmail(
		ctx,
		rawToken,
	)

	require.NoError(
		t,
		err,
	)

	// second verification must fail

	_, err = authService.VerifyEmail(
		ctx,
		rawToken,
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

func TestEmailVerificationExpired(t *testing.T) {

	ctx := context.Background()

	db := setupAuthTest(t)

	authService := helpers.NewAuthService(
		db,
	)

	req := service.RegisterRequest{
		Email: "verify-expired-" + uuid.NewString() + "@example.com",

		Username: "verifyexpired" + uuid.NewString()[:8],

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

	rawToken, verification :=
		helpers.CreateEmailVerification(
			user.ID,
		)

	// make token expired

	verification.ExpiresAt = time.Now().
		Add(
			-24 * time.Hour,
		)

	err = db.Create(
		verification,
	).Error

	require.NoError(
		t,
		err,
	)

	_, err = authService.VerifyEmail(
		ctx,
		rawToken,
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

func TestEmailVerificationTokenCannotBeUsedByAnotherUser(t *testing.T) {

	ctx := context.Background()

	db := setupAuthTest(t)

	authService := helpers.NewAuthService(
		db,
	)

	// User A

	userA, err := authService.Register(
		ctx,
		service.RegisterRequest{
			Email: "verify-a-" + uuid.NewString() + "@example.com",

			Username: "verifyusera" + uuid.NewString()[:8],

			Password: "C4ctus!River#829Lamp",
		},
	)

	require.NoError(
		t,
		err,
	)

	// User B

	userB, err := authService.Register(
		ctx,
		service.RegisterRequest{
			Email: "verify-b-" + uuid.NewString() + "@example.com",

			Username: "verifyuserb" + uuid.NewString()[:8],

			Password: "C4ctus!River#829Lamp",
		},
	)

	require.NoError(
		t,
		err,
	)

	// create token for User A

	rawToken, verification := helpers.CreateEmailVerification(
		userA.ID,
	)

	err = db.Create(
		verification,
	).Error

	require.NoError(
		t,
		err,
	)

	// ensure User B is still unverified

	require.Nil(
		t,
		userB.VerifiedAt,
	)

	// use User A token

	_, err = authService.VerifyEmail(
		ctx,
		rawToken,
	)

	require.NoError(
		t,
		err,
	)

	// reload users

	var updatedA models.User

	err = db.First(
		&updatedA,
		"id = ?",
		userA.ID,
	).Error

	require.NoError(
		t,
		err,
	)

	var updatedB models.User

	err = db.First(
		&updatedB,
		"id = ?",
		userB.ID,
	).Error

	require.NoError(
		t,
		err,
	)

	// A verified

	require.NotNil(
		t,
		updatedA.VerifiedAt,
	)

	// B must remain unverified

	assert.Nil(
		t,
		updatedB.VerifiedAt,
	)
}

func TestEmailVerificationAfterUserDeletion(t *testing.T) {

	ctx := context.Background()

	db := setupAuthTest(t)

	emailService := &helpers.CapturingEmailService{}

	authService := helpers.NewAuthServiceWithCapturingEmail(
		db,
		emailService,
	)

	userService := helpers.NewUserService(
		db,
	)

	// register user

	user, err := authService.Register(
		ctx,
		service.RegisterRequest{
			Email:    "verifydeleted@example.com",
			Username: "verifydeleted",
			Password: "C4ctus!River#829Lamp",
		},
	)

	require.NoError(
		t,
		err,
	)

	verificationToken := emailService.VerificationToken

	require.NotEmpty(
		t,
		verificationToken,
	)

	// delete user before verification

	err = userService.Delete(
		ctx,
		user.ID,
	)

	require.NoError(
		t,
		err,
	)

	// verification should fail

	_, err = authService.VerifyEmail(
		ctx,
		verificationToken,
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
