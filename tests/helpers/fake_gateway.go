package helpers

import (
	"context"
	"errors"

	"github.com/reinp/event-platform/backend/internal/payment/paypal"
)

type FakePaymentGateway struct {
	Order *paypal.Order

	CreateOrderCalled bool

	CapturedOrderID string

	VerifyWebhookCalled bool
}

func (f *FakePaymentGateway) CreateOrder(
	ctx context.Context,
	amount int64,
) (*paypal.Order, error) {

	f.CreateOrderCalled = true

	return f.Order, nil
}

func (f *FakePaymentGateway) CaptureOrder(
	ctx context.Context,
	orderID string,
) error {

	f.CapturedOrderID = orderID

	return nil
}

func (f *FakePaymentGateway) VerifyWebhookSignature(
	ctx context.Context,
	headers paypal.WebhookHeaders,
	body []byte,
) error {

	f.VerifyWebhookCalled = true

	return nil
}

type FailingPaymentGateway struct {
}

func (f *FailingPaymentGateway) CreateOrder(
	ctx context.Context,
	amount int64,
) (*paypal.Order, error) {

	return nil, errors.New("create order failed")
}

func (f *FailingPaymentGateway) CaptureOrder(
	ctx context.Context,
	orderID string,
) error {

	return errors.New("capture failed")
}

func (f *FailingPaymentGateway) VerifyWebhookSignature(
	ctx context.Context,
	headers paypal.WebhookHeaders,
	body []byte,
) error {

	return nil
}
