package paypal

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

func (c *Client) CaptureOrder(
	ctx context.Context,
	orderID string,
) (string, error) {

	token, err := c.GetAccessToken(ctx)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		c.baseURL+"/v2/checkout/orders/"+orderID+"/capture",
		nil,
	)
	if err != nil {
		return "", err
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
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusCreated {
		return "", fmt.Errorf(
			"paypal capture failed: status=%d body=%s",
			resp.StatusCode,
			string(body),
		)
	}

	var result struct {
		PurchaseUnits []struct {
			Payments struct {
				Captures []struct {
					ID string `json:"id"`
				} `json:"captures"`
			} `json:"payments"`
		} `json:"purchase_units"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return "", err
	}

	if len(result.PurchaseUnits) == 0 {
		return "", errors.New("paypal returned no purchase units")
	}

	if len(result.PurchaseUnits[0].Payments.Captures) == 0 {
		return "", errors.New("paypal returned no capture")
	}

	captureID := result.PurchaseUnits[0].Payments.Captures[0].ID

	return captureID, nil
}
