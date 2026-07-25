package paypal

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/reinp/event-platform/backend/internal/payment/paypal"
)

func TestPayPalWebhookVerifySignatureSuccess(t *testing.T) {

	server := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {

				switch r.URL.Path {

				case "/v1/oauth2/token":

					w.Header().Set(
						"Content-Type",
						"application/json",
					)

					w.WriteHeader(
						http.StatusOK,
					)

					w.Write([]byte(`
					{
						"access_token": "test-token",
						"token_type": "Bearer",
						"expires_in": 3600
					}
					`))

				case "/v1/notifications/verify-webhook-signature":

					w.Header().Set(
						"Content-Type",
						"application/json",
					)

					w.WriteHeader(
						http.StatusOK,
					)

					w.Write([]byte(`
					{
						"verification_status": "SUCCESS"
					}
					`))

				default:

					t.Fatalf(
						"unexpected path: %s",
						r.URL.Path,
					)
				}
			},
		),
	)

	defer server.Close()

	client := paypal.NewClientWithBaseURL(
		"test",
		"test",
		server.URL,
		"http://return",
		"http://cancel",
		"webhook-id",
	)

	headers := paypal.WebhookHeaders{

		TransmissionID: "test-transmission-id",

		TransmissionTime: "2026-01-01T00:00:00Z",

		CertURL: "https://example.com/cert",

		AuthAlgo: "SHA256withRSA",

		TransmissionSig: "test-signature",
	}

	body := []byte(`
	{
		"id": "WH-123",
		"event_type": "PAYMENT.CAPTURE.COMPLETED"
	}
	`)

	err := client.VerifyWebhookSignature(
		context.Background(),
		headers,
		body,
	)

	if err != nil {

		t.Fatalf(
			"expected webhook verification success, got: %v",
			err,
		)
	}
}

func TestPayPalWebhookVerifySignatureInvalidSignature(t *testing.T) {

	server := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {

				switch r.URL.Path {

				case "/v1/oauth2/token":

					w.Header().Set(
						"Content-Type",
						"application/json",
					)

					w.WriteHeader(
						http.StatusOK,
					)

					w.Write([]byte(`
					{
						"access_token": "test-token",
						"token_type": "Bearer",
						"expires_in": 3600
					}
					`))

				case "/v1/notifications/verify-webhook-signature":

					w.Header().Set(
						"Content-Type",
						"application/json",
					)

					w.WriteHeader(
						http.StatusOK,
					)

					w.Write([]byte(`
					{
						"verification_status": "FAILURE"
					}
					`))

				default:

					t.Fatalf(
						"unexpected path: %s",
						r.URL.Path,
					)
				}
			},
		),
	)

	defer server.Close()

	client := paypal.NewClientWithBaseURL(
		"test",
		"test",
		server.URL,
		"http://return",
		"http://cancel",
		"webhook-id",
	)

	headers := paypal.WebhookHeaders{

		TransmissionID: "invalid-transmission-id",

		TransmissionTime: "2026-01-01T00:00:00Z",

		CertURL: "https://example.com/cert",

		AuthAlgo: "SHA256withRSA",

		TransmissionSig: "invalid-signature",
	}

	body := []byte(`
	{
		"id": "WH-123",
		"event_type": "PAYMENT.CAPTURE.COMPLETED"
	}
	`)

	err := client.VerifyWebhookSignature(
		context.Background(),
		headers,
		body,
	)

	if err == nil {

		t.Fatal(
			"expected invalid signature error",
		)
	}
}

func TestPayPalWebhookVerifySignatureMissingHeaders(t *testing.T) {

	client := paypal.NewClientWithBaseURL(
		"test",
		"test",
		"http://unused",
		"http://return",
		"http://cancel",
		"webhook-id",
	)

	headers := paypal.WebhookHeaders{

		TransmissionID: "",

		TransmissionTime: "2026-01-01T00:00:00Z",

		CertURL: "https://example.com/cert",

		AuthAlgo: "SHA256withRSA",

		TransmissionSig: "signature",
	}

	body := []byte(`
	{
		"id": "WH-123",
		"event_type": "PAYMENT.CAPTURE.COMPLETED"
	}
	`)

	err := client.VerifyWebhookSignature(
		context.Background(),
		headers,
		body,
	)

	if err == nil {

		t.Fatal(
			"expected missing headers error",
		)
	}
}
