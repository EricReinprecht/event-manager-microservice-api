package paypal

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/reinp/event-platform/backend/internal/database"
	"github.com/reinp/event-platform/backend/internal/payment/paypal"
	"github.com/reinp/event-platform/backend/tests/helpers"
)

func TestPaymentCapturePaymentFailed(t *testing.T) {

	db, err := helpers.TestDatabaseSilent()

	if err != nil {
		t.Fatal(err)
	}

	executor := database.NewGormExecutor(db)

	purchaseService := helpers.NewPurchaseService(db)

	ticketService := helpers.NewTicketService(db)

	fakeGateway := &helpers.FailingPaymentGateway{}

	paymentService := helpers.NewPaymentService(
		executor,
		purchaseService,
		ticketService,
		fakeGateway,
	)

	_, err = paymentService.CapturePayment(
		context.Background(),
		"FAILED_ORDER_ID",
	)

	if err == nil {

		t.Fatal(
			"expected capture error",
		)
	}
}

func TestPayPalCaptureOrderAlreadyCaptured(t *testing.T) {

	server := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {

				switch r.URL.Path {

				case "/v1/oauth2/token":

					w.Header().Set(
						"Content-Type",
						"application/json",
					)

					w.WriteHeader(http.StatusOK)

					w.Write([]byte(`
					{
						"access_token": "test-token",
						"token_type": "Bearer",
						"expires_in": 3600
					}
					`))

				case "/v2/checkout/orders/ORDER123/capture":

					w.Header().Set(
						"Content-Type",
						"application/json",
					)

					w.WriteHeader(
						http.StatusUnprocessableEntity,
					)

					w.Write([]byte(`
					{
						"name": "UNPROCESSABLE_ENTITY",
						"details": [
							{
								"issue": "ORDER_ALREADY_CAPTURED"
							}
						]
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
		"webhook",
	)

	_, err := client.CaptureOrder(
		context.Background(),
		"ORDER123",
	)

	if err == nil {

		t.Fatal(
			"expected already captured error",
		)
	}
}

func TestPayPalCaptureOrderInvalidOrderID(t *testing.T) {

	server := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {

				switch r.URL.Path {

				case "/v1/oauth2/token":

					w.Header().Set(
						"Content-Type",
						"application/json",
					)

					w.WriteHeader(http.StatusOK)

					w.Write([]byte(`
					{
						"access_token": "test-token",
						"token_type": "Bearer",
						"expires_in": 3600
					}
					`))

				case "/v2/checkout/orders/INVALID/capture":

					w.Header().Set(
						"Content-Type",
						"application/json",
					)

					w.WriteHeader(
						http.StatusNotFound,
					)

					w.Write([]byte(`
					{
						"name": "RESOURCE_NOT_FOUND",
						"message": "The specified resource does not exist."
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
		"webhook",
	)

	_, err := client.CaptureOrder(
		context.Background(),
		"INVALID",
	)

	if err == nil {

		t.Fatal(
			"expected invalid order id error",
		)
	}
}

func TestPayPalCaptureOrderAPIError(t *testing.T) {

	server := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {

				switch r.URL.Path {

				case "/v1/oauth2/token":

					w.Header().Set(
						"Content-Type",
						"application/json",
					)

					w.WriteHeader(http.StatusOK)

					w.Write([]byte(`
					{
						"access_token": "test-token",
						"token_type": "Bearer",
						"expires_in": 3600
					}
					`))

				case "/v2/checkout/orders/ORDER123/capture":

					w.Header().Set(
						"Content-Type",
						"application/json",
					)

					w.WriteHeader(
						http.StatusInternalServerError,
					)

					w.Write([]byte(`
					{
						"name": "INTERNAL_SERVER_ERROR",
						"message": "Something went wrong"
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
		"webhook",
	)

	_, err := client.CaptureOrder(
		context.Background(),
		"ORDER123",
	)

	if err == nil {

		t.Fatal(
			"expected paypal api error",
		)
	}
}
