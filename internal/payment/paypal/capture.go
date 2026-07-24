package paypal

import (
	"context"
	"errors"
	"net/http"
)

func (c *Client) CaptureOrder(
	ctx context.Context,
	orderID string,
) error {

	token, err := c.GetAccessToken(ctx)

	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		c.baseURL+"/v2/checkout/orders/"+orderID+"/capture",
		nil,
	)

	if err != nil {
		return err
	}

	req.Header.Set(
		"Authorization",
		"Bearer "+token,
	)

	req.Header.Set(
		"Content-Type",
		"application/json",
	)

	req.Header.Set(
		"Prefer",
		"return=representation",
	)

	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {

		return errors.New(
			"paypal capture failed",
		)
	}

	return nil
}
