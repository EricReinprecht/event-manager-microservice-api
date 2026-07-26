package integration

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"

	"github.com/gin-gonic/gin"
	"github.com/reinp/event-platform/backend/internal/middleware"
	"github.com/reinp/event-platform/backend/tests/helpers"
)

func TestUnauthenticatedRequestBlockedByAuthMiddleware(t *testing.T) {

	gin.SetMode(gin.TestMode)

	db, err := helpers.TestDatabaseSilent()

	if err != nil {
		t.Fatal(err)
	}

	authService := helpers.NewAuthService(
		db,
	)

	router := gin.New()

	router.Use(
		middleware.Auth(
			authService,
		),
	)

	router.POST(
		"/api/purchases/:id/refund",
		func(c *gin.Context) {

			t.Fatal(
				"handler should not be reached",
			)

		},
	)

	req := httptest.NewRequest(
		http.MethodPost,
		"/api/purchases/"+uuid.New().String()+"/refund",
		nil,
	)

	recorder := httptest.NewRecorder()

	router.ServeHTTP(
		recorder,
		req,
	)

	if recorder.Code != http.StatusUnauthorized {

		t.Fatalf(
			"expected 401 got %d body=%s",
			recorder.Code,
			recorder.Body.String(),
		)
	}
}
