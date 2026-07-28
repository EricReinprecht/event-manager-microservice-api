package auth_test

import (
	"fmt"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/reinp/event-platform/backend/tests/helpers"
	"github.com/reinp/event-platform/backend/tests/helpers/http"
)

func TestUnauthorizedRequests(t *testing.T) {

	gin.SetMode(gin.TestMode)

	db := setupAuthTest(t)

	authService := helpers.NewAuthService(
		db,
	)

	router := http.ProtectedAuthRouter(
		authService,
	)

	tests := []struct {
		name   string
		method string
		url    string
	}{
		{
			name:   "get sessions without token",
			method: "GET",
			url:    "/api/auth/sessions",
		},
		{
			name:   "logout without token",
			method: "DELETE",
			url:    "/api/auth/sessions",
		},
	}

	for _, tt := range tests {

		t.Run(
			tt.name,
			func(t *testing.T) {

				req := helpers.JSONRequest(
					tt.method,
					tt.url,
					nil,
				)

				resp := helpers.ExecuteRequest(
					router,
					req,
				)

				require.Equal(
					t,
					401,
					resp.Code,
				)

				assert.Contains(
					t,
					resp.Body.String(),
					"error",
				)
			},
		)
	}
}

func TestUnauthorizedInvalidToken(t *testing.T) {

	gin.SetMode(gin.TestMode)

	db := setupAuthTest(t)

	authService := helpers.NewAuthService(
		db,
	)

	router := http.ProtectedAuthRouter(
		authService,
	)

	resp := helpers.DoAuthenticatedRequest(
		t,
		router,
		"GET",
		"/api/auth/sessions",
		"invalid.jwt.token",
		nil,
	)

	require.Equal(
		t,
		401,
		resp.StatusCode,
	)

	assert.Contains(
		t,
		fmt.Sprint(resp.Body),
		"invalid token",
	)
}
