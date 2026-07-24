package paypal

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
)

type Client struct {
	clientID     string
	clientSecret string
	baseURL      string

	returnURL string
	cancelURL string
	webhookID string
}

func NewClient(
	clientID string,
	clientSecret string,
	baseURL string,
	returnURL string,
	cancelURL string,
	webhookID string,
) *Client {

	return &Client{
		clientID:     clientID,
		clientSecret: clientSecret,
		baseURL:      baseURL,
		returnURL:    returnURL,
		cancelURL:    cancelURL,
		webhookID:    webhookID,
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

	token, err := c.GetAccessToken(ctx)

	if err != nil {
		return nil, err
	}

	value := strconv.FormatFloat(
		float64(amount)/100,
		'f',
		2,
		64,
	)

	payload := map[string]interface{}{

		"intent": "CAPTURE",

		"purchase_units": []interface{}{

			map[string]interface{}{

				"amount": map[string]interface{}{

					"currency_code": "EUR",

					"value": value,
				},
			},
		},

		"application_context": map[string]interface{}{

			"return_url": c.returnURL,

			"cancel_url": c.cancelURL,

			"user_action": "PAY_NOW",
		},
	}

	body, err := json.Marshal(payload)

	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		c.baseURL+"/v2/checkout/orders",
		bytes.NewBuffer(body),
	)

	if err != nil {
		return nil, err
	}

	req.Header.Set(
		"Authorization",
		"Bearer "+token,
	)

	req.Header.Set(
		"Content-Type",
		"application/json",
	)

	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {

		return nil,
			errors.New(
				"paypal order creation failed",
			)
	}

	var result struct {
		ID string `json:"id"`

		Links []struct {
			Href string `json:"href"`

			Rel string `json:"rel"`
		} `json:"links"`
	}

	if err := json.NewDecoder(
		resp.Body,
	).Decode(&result); err != nil {

		return nil, err
	}

	order := &Order{
		ID: result.ID,
	}

	for _, link := range result.Links {

		if link.Rel == "approve" {

			order.ApprovalURL = link.Href
			break
		}
	}

	if order.ID == "" {

		return nil,
			errors.New(
				"paypal returned empty order id",
			)
	}

	return order, nil
}

func (c *Client) VerifyWebhookSignature(
	ctx context.Context,
	headers WebhookHeaders,
	body []byte,
) error {

	accessToken, err := c.GetAccessToken(ctx)

	if err != nil {
		return err
	}

	payload := map[string]interface{}{
		"auth_algo": headers.AuthAlgo,

		"cert_url": headers.CertURL,

		"transmission_id": headers.TransmissionID,

		"transmission_sig": headers.TransmissionSig,

		"transmission_time": headers.TransmissionTime,

		"webhook_id": c.webhookID,

		"webhook_event": json.RawMessage(body),
	}

	data, err := json.Marshal(payload)

	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		c.baseURL+"/v1/notifications/verify-webhook-signature",
		bytes.NewBuffer(data),
	)

	if err != nil {
		return err
	}

	req.Header.Set(
		"Authorization",
		"Bearer "+accessToken,
	)

	req.Header.Set(
		"Content-Type",
		"application/json",
	)

	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {

		return errors.New(
			"paypal webhook verification request failed",
		)
	}

	var result struct {
		VerificationStatus string `json:"verification_status"`
	}

	err = json.NewDecoder(
		resp.Body,
	).Decode(&result)

	if err != nil {
		return err
	}

	if result.VerificationStatus != "SUCCESS" {

		return errors.New(
			"paypal webhook signature invalid",
		)
	}

	return nil
}
