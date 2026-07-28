package auth_test

import (
	"context"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/reinp/event-platform/backend/internal/models"
	"github.com/reinp/event-platform/backend/tests/helpers"
	"github.com/reinp/event-platform/backend/tests/helpers/http"
	"github.com/reinp/event-platform/backend/tests/scenarios"
)

func TestGetSessions(
	t *testing.T,
) {

	gin.SetMode(gin.TestMode)

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

	var refreshToken models.RefreshToken

	err := db.
		Where(
			"user_id = ?",
			user.ID,
		).
		Order(
			"created_at DESC",
		).
		First(
			&refreshToken,
		).
		Error

	require.NoError(
		t,
		err,
	)

	sessionHandler := helpers.NewSessionHandler(
		db,
	)

	router := http.AuthRouter(
		sessionHandler,
		user.ID,
		refreshToken.FamilyID,
	)

	resp := helpers.DoAuthenticatedRequest(
		t,
		router,
		"GET",
		"/api/auth/sessions",
		login.AccessToken,
		nil,
	)

	require.Equal(
		t,
		200,
		resp.StatusCode,
	)

	assert.NotEmpty(
		t,
		resp.Body,
	)
}

func TestLogoutAllSessions(
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

	sessionHandler := helpers.NewSessionHandler(
		db,
	)

	router := http.AuthRouter(
		sessionHandler,
		user.ID,
		uuid.Nil,
	)

	resp := helpers.DoAuthenticatedRequest(
		t,
		router,
		"DELETE",
		"/api/auth/sessions",
		login.AccessToken,
		nil,
	)

	require.Equal(
		t,
		200,
		resp.StatusCode,
	)

	_, err := authService.Refresh(
		ctx,
		login.RefreshToken,
	)

	require.Error(
		t,
		err,
	)
}
