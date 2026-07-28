package auth_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/reinp/event-platform/backend/internal/service"
	"github.com/reinp/event-platform/backend/tests/helpers"
	"github.com/reinp/event-platform/backend/tests/scenarios"
)

func TestLoginSuccess(t *testing.T) {

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

func TestLoginWrongPassword(t *testing.T) {

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

	_, err := authService.Login(
		ctx,
		user.Email,
		"WrongPassword123!",
		"test-agent",
		"127.0.0.1",
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

func TestLoginUnverifiedUser(t *testing.T) {

	ctx := context.Background()

	db := setupAuthTest(t)

	authService := helpers.NewAuthService(
		db,
	)

	// create user through registration
	// (Register creates an unverified user)

	req := service.RegisterRequest{
		Email: "unverified-" + uuid.NewString() + "@example.com",

		Username: "unverified" + uuid.NewString()[:8],

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

	require.Nil(
		t,
		user.VerifiedAt,
	)

	_, err = authService.Login(
		ctx,
		user.Email,
		req.Password,
		"test-agent",
		"127.0.0.1",
	)

	require.Error(
		t,
		err,
	)

	assert.Contains(
		t,
		err.Error(),
		"verified",
	)
}

func TestDeletedUserCannotLogin(t *testing.T) {

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

	// login should fail

	_, err = authService.Login(
		ctx,
		user.Email,
		"C4ctus!River#829Lamp",
		"test-agent",
		"127.0.0.1",
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

func TestLoginDeletedUserAfterRefreshTokenExists(t *testing.T) {

	ctx := context.Background()

	db := setupAuthTest(t)

	authService := helpers.NewAuthService(
		db,
	)

	// create user
	req := service.RegisterRequest{
		Email: "deleted-login-" + uuid.NewString() + "@example.com",

		Username: "deletedlogin" + uuid.NewString()[:8],

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
	user.VerifiedAt = helpers.Ptr(time.Now())

	err = db.Save(user).Error

	require.NoError(
		t,
		err,
	)

	// login once -> creates refresh token
	tokens, err := authService.Login(
		ctx,
		user.Email,
		req.Password,
		"test-agent",
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

	// delete user
	err = db.Delete(
		user,
	).Error

	require.NoError(
		t,
		err,
	)

	// login should fail
	_, err = authService.Login(
		ctx,
		req.Email,
		req.Password,
		"test-agent",
		"127.0.0.1",
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

func TestLoginCaseInsensitiveEmail(t *testing.T) {

	ctx := context.Background()

	db := setupAuthTest(t)

	authService := helpers.NewAuthService(
		db,
	)

	email := "CaseSensitiveUser@example.com"

	user, err := authService.Register(
		ctx,
		service.RegisterRequest{
			Email: email,

			Username: "caseemailuser" + uuid.NewString()[:8],

			Password: "C4ctus!River#829Lamp",
		},
	)

	require.NoError(
		t,
		err,
	)

	// verify email

	now := time.Now()

	user.VerifiedAt = &now

	err = db.Save(
		user,
	).Error

	require.NoError(
		t,
		err,
	)

	// login with different email casing

	response, err := authService.Login(
		ctx,
		"casesensitiveUSER@EXAMPLE.COM",
		"C4ctus!River#829Lamp",
		"test-browser",
		"127.0.0.1",
	)

	require.NoError(
		t,
		err,
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
}
