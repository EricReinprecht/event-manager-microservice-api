package auth_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/reinp/event-platform/backend/internal/service"

	"github.com/reinp/event-platform/backend/tests/helpers"
)

func TestRegisterSuccess(
	t *testing.T,
) {

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

func TestRegisterDuplicateEmail(
	t *testing.T,
) {

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

func TestRegisterDuplicateUsername(
	t *testing.T,
) {

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

func TestRegisterWeakPassword(
	t *testing.T,
) {

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

func TestRegisterCreatesUnverifiedUser(
	t *testing.T,
) {

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
