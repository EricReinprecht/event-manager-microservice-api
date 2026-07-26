package payment

import (
	"context"

	"github.com/reinp/event-platform/backend/internal/payment/paypal"
)

type Gateway interface {
	CreateOrder(
		ctx context.Context,
		amount int64,
	) (*paypal.Order, error)

	CaptureOrder(
		ctx context.Context,
		orderID string,
	) (string, error)

	VerifyWebhookSignature(
		ctx context.Context,
		headers paypal.WebhookHeaders,
		body []byte,
	) error

	RefundPayment(
		ctx context.Context,
		paymentID string,
	) (string, error)
}
