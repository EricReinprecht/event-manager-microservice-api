package helpers

import (
	"os"

	"github.com/reinp/event-platform/backend/internal/payment/paypal"
)

func NewPayPalClient() *paypal.Client {

	return paypal.NewClient(
		os.Getenv("PAYPAL_CLIENT_ID"),
		os.Getenv("PAYPAL_CLIENT_SECRET"),
		os.Getenv("PAYPAL_BASE_URL"),
	)
}
