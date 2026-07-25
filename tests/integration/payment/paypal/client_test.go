package paypal

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/reinp/event-platform/backend/internal/payment/paypal"
)

func TestPayPalGetAccessTokenSuccess(t *testing.T) {

	server := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {

				if r.URL.Path != "/v1/oauth2/token" {

					t.Fatalf(
						"unexpected path: %s",
						r.URL.Path,
					)
				}

				if r.Method != http.MethodPost {

					t.Fatalf(
						"expected POST request, got %s",
						r.Method,
					)
				}

				w.Header().
					Set(
						"Content-Type",
						"application/json",
					)

				json.NewEncoder(w).
					Encode(
						map[string]interface{}{
							"access_token": "test-access-token",
							"token_type":   "Bearer",
							"expires_in":   3600,
						},
					)
			},
		),
	)

	defer server.Close()

	client := paypal.NewClient(
		"client-id",
		"client-secret",
		server.URL,
		"",
		"",
		"",
	)

	token, err := client.GetAccessToken(
		context.Background(),
	)

	if err != nil {
		t.Fatal(err)
	}

	if token != "test-access-token" {

		t.Fatalf(
			"expected access token test-access-token, got %s",
			token,
		)
	}
}

func TestPayPalGetAccessTokenInvalidCredentials(t *testing.T) {

	server := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {

				if r.URL.Path != "/v1/oauth2/token" {

					t.Fatalf(
						"unexpected path: %s",
						r.URL.Path,
					)
				}

				if r.Method != http.MethodPost {

					t.Fatalf(
						"expected POST request, got %s",
						r.Method,
					)
				}

				w.Header().
					Set(
						"Content-Type",
						"application/json",
					)

				w.WriteHeader(
					http.StatusUnauthorized,
				)

				json.NewEncoder(w).
					Encode(
						map[string]interface{}{
							"name":    "AUTHENTICATION_FAILURE",
							"message": "Authentication failed due to invalid credentials",
						},
					)
			},
		),
	)

	defer server.Close()

	client := paypal.NewClient(
		"invalid-client-id",
		"invalid-client-secret",
		server.URL,
		"",
		"",
		"",
	)

	_, err := client.GetAccessToken(
		context.Background(),
	)

	if err == nil {

		t.Fatal(
			"expected authentication error",
		)
	}
}

func TestPayPalGetAccessTokenServerError(t *testing.T) {

	server := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {

				if r.URL.Path != "/v1/oauth2/token" {

					t.Fatalf(
						"unexpected path: %s",
						r.URL.Path,
					)
				}

				if r.Method != http.MethodPost {

					t.Fatalf(
						"expected POST request, got %s",
						r.Method,
					)
				}

				w.Header().
					Set(
						"Content-Type",
						"application/json",
					)

				w.WriteHeader(
					http.StatusInternalServerError,
				)

				json.NewEncoder(w).
					Encode(
						map[string]interface{}{
							"name":    "INTERNAL_SERVER_ERROR",
							"message": "PayPal server error",
						},
					)
			},
		),
	)

	defer server.Close()

	client := paypal.NewClient(
		"client-id",
		"client-secret",
		server.URL,
		"",
		"",
		"",
	)

	_, err := client.GetAccessToken(
		context.Background(),
	)

	if err == nil {

		t.Fatal(
			"expected server error from PayPal",
		)
	}
}
