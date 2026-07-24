package paypal

import "context"

type Client struct {
	clientID     string
	clientSecret string
	baseURL      string
}

func NewClient(
	clientID string,
	clientSecret string,
	baseURL string,
) *Client {

	return &Client{
		clientID:     clientID,
		clientSecret: clientSecret,
		baseURL:      baseURL,
	}
}

type Order struct {
	ID string

	ApprovalURL string
}

func (c *Client) CreateOrder(
	ctx context.Context,
	amount int64,
) (*Order, error) {

	// TODO:
	// call PayPal API

	return &Order{
		ID: "PAYPAL_TEST_ORDER_ID",

		ApprovalURL: "https://paypal.test/checkout",
	}, nil
}

func (c *Client) VerifyWebhookSignature(
	ctx context.Context,
	headers WebhookHeaders,
	body []byte,
) error {

	// call PayPal verification endpoint here

	return nil
}
