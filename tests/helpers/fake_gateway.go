package helpers

import (
	"context"
	"errors"

	"github.com/google/uuid"

	"github.com/reinp/event-platform/backend/internal/database"
	"github.com/reinp/event-platform/backend/internal/payment/paypal"
)

type FakePaymentGateway struct {
	Order *paypal.Order

	CreateOrderCalled bool

	CaptureOrderCalled bool

	CapturedOrderID string

	RefundCalled bool

	RefundedPaymentID string

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
) (string, error) {

	f.CaptureOrderCalled = true

	f.CapturedOrderID = orderID

	return "FAKE-CAPTURE-ID", nil
}

func (f *FakePaymentGateway) RefundPayment(
	ctx context.Context,
	paymentID string,
	amount int64,
) (string, error) {

	f.RefundCalled = true

	f.RefundedPaymentID = paymentID

	return "FAKE-REFUND-ID", nil
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

	return nil, errors.New(
		"create order failed",
	)
}

func (f *FailingPaymentGateway) CaptureOrder(
	ctx context.Context,
	orderID string,
) (string, error) {

	return "", errors.New(
		"capture failed",
	)
}

func (f *FailingPaymentGateway) RefundPayment(
	ctx context.Context,
	paymentID string,
	amount int64,
) (string, error) {

	return "", errors.New(
		"refund failed",
	)
}

func (f *FailingPaymentGateway) VerifyWebhookSignature(
	ctx context.Context,
	headers paypal.WebhookHeaders,
	body []byte,
) error {

	return nil
}

// Used for webhook signature rejection tests

type InvalidSignaturePaymentGateway struct {
}

func (f *InvalidSignaturePaymentGateway) CreateOrder(
	ctx context.Context,
	amount int64,
) (*paypal.Order, error) {

	return nil, errors.New(
		"not implemented",
	)
}

func (f *InvalidSignaturePaymentGateway) CaptureOrder(
	ctx context.Context,
	orderID string,
) (string, error) {

	return "", errors.New(
		"not implemented",
	)
}

func (f *InvalidSignaturePaymentGateway) RefundPayment(
	ctx context.Context,
	paymentID string,
	amount int64,
) (string, error) {

	return "", errors.New(
		"not implemented",
	)
}

func (f *InvalidSignaturePaymentGateway) VerifyWebhookSignature(
	ctx context.Context,
	headers paypal.WebhookHeaders,
	body []byte,
) error {

	return errors.New(
		"invalid paypal signature",
	)
}

type FailingTicketRepository struct {
}

func (r *FailingTicketRepository) CancelByPurchase(
	tx database.DBExecutor,
	purchaseID uuid.UUID,
) error {

	return errors.New(
		"ticket cancellation failed",
	)
}
