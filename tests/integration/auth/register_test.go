package auth_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/reinp/event-platform/backend/internal/models"
	"github.com/reinp/event-platform/backend/internal/service"

	"github.com/reinp/event-platform/backend/tests/helpers"
)

func TestRegisterSuccess(t *testing.T) {

	ctx := context.Background()

	db := setupAuthTest(t)

	authService := helpers.NewAuthService(
		db,
	)

	req := service.RegisterRequest{
		Email: "new-" + uuid.NewString() + "@example.com",

		Username: "newuser" + uuid.NewString()[:8],

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

	require.NotNil(
		t,
		user,
	)

	assert.NotEqual(
		t,
		uuid.Nil,
		user.ID,
	)

	assert.Equal(
		t,
		req.Email,
		user.Email,
	)

	assert.Equal(
		t,
		req.Username,
		user.Username,
	)

	assert.Nil(
		t,
		user.VerifiedAt,
	)
}

func TestRegisterDuplicateEmail(t *testing.T) {

	ctx := context.Background()

	db := setupAuthTest(t)

	authService := helpers.NewAuthService(
		db,
	)

	first, err := authService.Register(
		ctx,
		service.RegisterRequest{
			Email: "first@example.com",

			Username: "firstuser",

			Password: "C4ctus!River#829Lamp",
		},
	)

	require.NoError(
		t,
		err,
	)

	_, err = authService.Register(
		ctx,
		service.RegisterRequest{
			Email: first.Email,

			Username: "anotheruser",

			Password: "C4ctus!River#829Lamp",
		},
	)

	require.Error(
		t,
		err,
	)

	assert.Contains(
		t,
		err.Error(),
		"email",
	)
}

func TestRegisterDuplicateUsername(t *testing.T) {

	ctx := context.Background()

	db := setupAuthTest(t)

	authService := helpers.NewAuthService(
		db,
	)

	first, err := authService.Register(
		ctx,
		service.RegisterRequest{
			Email: "first@example.com",

			Username: "firstuser",

			Password: "C4ctus!River#829Lamp",
		},
	)

	require.NoError(
		t,
		err,
	)

	_, err = authService.Register(
		ctx,
		service.RegisterRequest{
			Email: "unique@example.com",

			Username: first.Username,

			Password: "C4ctus!River#829Lamp",
		},
	)

	require.Error(
		t,
		err,
	)

	assert.Contains(
		t,
		err.Error(),
		"username",
	)
}

func TestRegisterWeakPassword(t *testing.T) {

	ctx := context.Background()

	db := setupAuthTest(t)

	authService := helpers.NewAuthService(
		db,
	)

	_, err := authService.Register(
		ctx,
		service.RegisterRequest{
			Email: "weak@example.com",

			Username: "weakuser",

			Password: "123",
		},
	)

	require.Error(
		t,
		err,
	)
}

func TestRegisterCreatesUnverifiedUser(t *testing.T) {

	ctx := context.Background()

	db := setupAuthTest(t)

	authService := helpers.NewAuthService(
		db,
	)

	user, err := authService.Register(
		ctx,
		service.RegisterRequest{
			Email: "verify@example.com",

			Username: "verifyuser",

			Password: "C4ctus!River#829Lamp",
		},
	)

	require.NoError(
		t,
		err,
	)

	assert.Nil(
		t,
		user.VerifiedAt,
	)

	var count int64

	err = db.
		Table(
			"email_verifications",
		).
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

func TestRegisterFailsWhenMailFails(t *testing.T) {

	ctx := context.Background()

	db := setupAuthTest(t)

	authService := helpers.NewAuthServiceWithEmailService(
		db,
		&helpers.FailingEmailService{},
	)

	_, err := authService.Register(
		ctx,
		service.RegisterRequest{
			Email: "mailfail@example.com",

			Username: "mailfailuser",

			Password: "C4ctus!River#829Lamp",
		},
	)

	require.Error(
		t,
		err,
	)

	assert.Contains(
		t,
		err.Error(),
		"mail",
	)

	// user should not remain created

	var count int64

	err = db.
		Table("users").
		Where(
			"email = ?",
			"mailfail@example.com",
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
		int64(0),
		count,
	)
}

func TestRegisterInvalidEmail(t *testing.T) {

	ctx := context.Background()

	db := setupAuthTest(t)

	authService := helpers.NewAuthService(
		db,
	)

	email := "invalid-email"

	username := "invalidemailuser"

	_, err := authService.Register(
		ctx,
		service.RegisterRequest{
			Email: email,

			Username: username,

			Password: "C4ctus!River#829Lamp",
		},
	)

	require.Error(
		t,
		err,
	)

	// make sure user was not created

	var userCount int64

	err = db.
		Table("users").
		Where(
			"email = ?",
			email,
		).
		Count(
			&userCount,
		).
		Error

	require.NoError(
		t,
		err,
	)

	assert.Equal(
		t,
		int64(0),
		userCount,
	)
}

func TestRegisterInvalidUsername(t *testing.T) {

	ctx := context.Background()

	db := setupAuthTest(t)

	authService := helpers.NewAuthService(
		db,
	)

	email := "invalid-username-" + uuid.NewString() + "@example.com"

	username := "ab" // less than minimum 3 chars

	_, err := authService.Register(
		ctx,
		service.RegisterRequest{
			Email: email,

			Username: username,

			Password: "C4ctus!River#829Lamp",
		},
	)

	require.Error(
		t,
		err,
	)

	assert.Contains(
		t,
		err.Error(),
		"username",
	)

	// ensure user was not created

	var count int64

	err = db.
		Table("users").
		Where(
			"email = ?",
			email,
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
		int64(0),
		count,
	)
}

func TestRegisterEmailNormalization(t *testing.T) {

	ctx := context.Background()

	db := setupAuthTest(t)

	authService := helpers.NewAuthService(
		db,
	)

	email := "  TestUser@Example.COM  "

	user, err := authService.Register(
		ctx,
		service.RegisterRequest{
			Email: email,

			Username: "emailnormaluser" + uuid.NewString()[:8],

			Password: "C4ctus!River#829Lamp",
		},
	)

	require.NoError(
		t,
		err,
	)

	require.NotNil(
		t,
		user,
	)

	// reload from database

	var stored models.User

	err = db.First(
		&stored,
		"id = ?",
		user.ID,
	).Error

	require.NoError(
		t,
		err,
	)

	assert.Equal(
		t,
		"testuser@example.com",
		stored.Email,
	)
}

func TestRegisterUsernameNormalization(t *testing.T) {

	ctx := context.Background()

	db := setupAuthTest(t)

	authService := helpers.NewAuthService(
		db,
	)

	username := "  normalized_user  "

	user, err := authService.Register(
		ctx,
		service.RegisterRequest{
			Email: "username-normalize-" + uuid.NewString() + "@example.com",

			Username: username,

			Password: "C4ctus!River#829Lamp",
		},
	)

	require.NoError(
		t,
		err,
	)

	require.NotNil(
		t,
		user,
	)

	// reload from database

	var stored models.User

	err = db.First(
		&stored,
		"id = ?",
		user.ID,
	).Error

	require.NoError(
		t,
		err,
	)

	assert.Equal(
		t,
		"normalized_user",
		stored.Username,
	)
}
