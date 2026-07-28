package auth_test

import (
	"context"
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/reinp/event-platform/backend/internal/auth"
	"github.com/reinp/event-platform/backend/internal/middleware"
	"github.com/reinp/event-platform/backend/internal/models"
	"github.com/reinp/event-platform/backend/internal/service"
	"github.com/reinp/event-platform/backend/tests/helpers"
	"github.com/reinp/event-platform/backend/tests/helpers/http"
	"github.com/reinp/event-platform/backend/tests/scenarios"
)

func TestGetSessions(t *testing.T) {

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

func TestLogoutAllSessions(t *testing.T) {

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

func TestLogoutSingleSession(t *testing.T) {

	gin.SetMode(gin.TestMode)

	ctx := context.Background()

	db := setupAuthTest(t)

	authService := helpers.NewAuthService(
		db,
	)

	sessionHandler := helpers.NewSessionHandler(
		db,
	)

	scenario := scenarios.NewAuthScenario(
		db,
		authService,
	)

	user := scenario.CreateVerifiedUser(
		ctx,
	)

	// create first session
	loginOne := scenario.Login(
		ctx,
		user,
	)

	// create second session
	loginTwo := scenario.Login(
		ctx,
		user,
	)

	// find first session family
	var refreshToken models.RefreshToken

	err := db.
		Where(
			"user_id = ?",
			user.ID,
		).
		Order(
			"created_at ASC",
		).
		First(
			&refreshToken,
		).
		Error

	require.NoError(
		t,
		err,
	)

	familyID := refreshToken.FamilyID

	require.NoError(
		t,
		err,
	)

	router := http.AuthRouter(
		sessionHandler,
		user.ID,
		familyID,
	)

	resp := helpers.DoAuthenticatedRequest(
		t,
		router,
		"DELETE",
		"/api/auth/sessions/"+familyID.String(),
		loginOne.AccessToken,
		nil,
	)

	require.Equal(
		t,
		200,
		resp.StatusCode,
	)

	// revoked session refresh token should fail

	_, err = authService.Refresh(
		ctx,
		loginOne.RefreshToken,
	)

	require.Error(
		t,
		err,
	)

	// other session should still work

	_, err = authService.Refresh(
		ctx,
		loginTwo.RefreshToken,
	)

	require.NoError(
		t,
		err,
	)
}

func TestCannotDeleteOtherUsersSession(t *testing.T) {

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

	// User A

	userA := scenario.CreateVerifiedUser(
		ctx,
	)

	// User B

	userB := scenario.CreateVerifiedUser(
		ctx,
	)

	loginB := scenario.Login(
		ctx,
		userB,
	)

	var session struct {
		FamilyID uuid.UUID
	}

	err := db.
		Table("refresh_tokens").
		Select("family_id").
		Where(
			"user_id = ?",
			userB.ID,
		).
		Order(
			"created_at DESC",
		).
		First(
			&session,
		).
		Error

	require.NoError(
		t,
		err,
	)

	userBFamilyID := session.FamilyID

	require.NoError(
		t,
		err,
	)

	sessionHandler := helpers.NewSessionHandler(
		db,
	)

	router := http.AuthRouter(
		sessionHandler,
		userA.ID,
		userA.ID,
	)

	resp := helpers.DoAuthenticatedRequest(
		t,
		router,
		"DELETE",
		"/api/auth/sessions/"+userBFamilyID.String(),
		loginB.AccessToken,
		nil,
	)

	require.NotEqual(
		t,
		200,
		resp.StatusCode,
	)
}

func TestCannotViewOtherUsersSessions(t *testing.T) {

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

	// User A

	userA := scenario.CreateVerifiedUser(
		ctx,
	)

	loginA := scenario.Login(
		ctx,
		userA,
	)

	// User B

	userB := scenario.CreateVerifiedUser(
		ctx,
	)

	scenario.Login(
		ctx,
		userB,
	)

	// Create real session handler

	sessionHandler := helpers.NewSessionHandler(
		db,
	)

	router := gin.New()

	router.GET(
		"/api/auth/sessions",
		middleware.Auth(authService),
		sessionHandler.GetSessions,
	)

	// User A requests sessions
	// The middleware identifies User A from the token

	resp := helpers.DoAuthenticatedRequest(
		t,
		router,
		"GET",
		"/api/auth/sessions",
		loginA.AccessToken,
		nil,
	)

	require.Equal(
		t,
		200,
		resp.StatusCode,
	)

	// decode response

	sessions, ok := resp.Body.([]interface{})

	require.True(
		t,
		ok,
	)
	// User A should not see User B sessions

	for _, item := range sessions {

		session, ok := item.(map[string]interface{})

		require.True(
			t,
			ok,
		)

		assert.NotEqual(
			t,
			userB.ID.String(),
			session["userId"],
		)
	}
}

func TestSessionListDoesNotLeakTokens(t *testing.T) {

	ctx := context.Background()

	db := setupAuthTest(t)

	authService := helpers.NewAuthService(
		db,
	)

	user, err := authService.Register(
		ctx,
		service.RegisterRequest{
			Email: "session-leak-" + uuid.NewString() + "@example.com",

			Username: "sessionleak" + uuid.NewString()[:8],

			Password: "C4ctus!River#829Lamp",
		},
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

	// create session

	tokens, err := authService.Login(
		ctx,
		user.Email,
		"C4ctus!River#829Lamp",
		"test-browser",
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

	// get sessions

	sessions, err := authService.Sessions(
		ctx,
		user.ID,
		uuid.Nil,
	)

	require.NoError(
		t,
		err,
	)

	require.Len(
		t,
		sessions,
		1,
	)

	session := sessions[0]

	// verify expected session data exists

	assert.Equal(
		t,
		"test-browser",
		session.Device,
	)

	assert.Equal(
		t,
		"127.0.0.1",
		session.IP,
	)

	// serialize response

	jsonData, err := json.Marshal(
		session,
	)

	require.NoError(
		t,
		err,
	)

	jsonString := string(jsonData)

	// security checks:
	// never expose tokens or hashes

	assert.NotContains(
		t,
		jsonString,
		"refresh",
	)

	assert.NotContains(
		t,
		jsonString,
		"token",
	)

	assert.NotContains(
		t,
		jsonString,
		"hash",
	)
}

func TestLogoutInvalidToken(t *testing.T) {

	ctx := context.Background()

	db := setupAuthTest(t)

	authService := helpers.NewAuthService(
		db,
	)

	_, err := authService.Register(
		ctx,
		service.RegisterRequest{
			Email: "logout-invalid-" + uuid.NewString() + "@example.com",

			Username: "logoutinvalid" + uuid.NewString()[:8],

			Password: "C4ctus!River#829Lamp",
		},
	)

	require.NoError(
		t,
		err,
	)

	// create random invalid refresh token

	invalidToken := auth.GenerateToken()

	err = authService.Logout(
		ctx,
		invalidToken,
	)

	require.Error(
		t,
		err,
	)

	assert.Contains(
		t,
		err.Error(),
		"invalid refresh token",
	)
}

func TestLogoutDoesNotAffectOtherSessions(t *testing.T) {

	ctx := context.Background()

	db := setupAuthTest(t)

	authService := helpers.NewAuthService(
		db,
	)

	user, err := authService.Register(
		ctx,
		service.RegisterRequest{
			Email: "logout-isolation-" + uuid.NewString() + "@example.com",

			Username: "logoutuser" + uuid.NewString()[:8],

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

	// create session A

	sessionA, err := authService.Login(
		ctx,
		user.Email,
		"C4ctus!River#829Lamp",
		"device-a",
		"127.0.0.1",
	)

	require.NoError(
		t,
		err,
	)

	// create session B

	sessionB, err := authService.Login(
		ctx,
		user.Email,
		"C4ctus!River#829Lamp",
		"device-b",
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

	// logout only session A

	err = authService.Logout(
		ctx,
		sessionA.RefreshToken,
	)

	require.NoError(
		t,
		err,
	)

	// session A should fail

	_, err = authService.Refresh(
		ctx,
		sessionA.RefreshToken,
	)

	require.Error(
		t,
		err,
	)

	assert.True(
		t,
		strings.Contains(err.Error(), "invalid") ||
			strings.Contains(err.Error(), "refresh token replay detected"),
	)

	// session B should still work

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

	assert.NotEmpty(
		t,
		newTokens.RefreshToken,
	)
}
