package paypal

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strings"
)

type accessTokenResponse struct {
	AccessToken string `json:"accessToken"`
	TokenType   string `json:"tokenType"`
	ExpiresIn   int    `json:"expiresIn"`
}

func (c *Client) GetAccessToken(
	ctx context.Context,
) (string, error) {

	data := url.Values{}

	data.Set(
		"grant_type",
		"client_credentials",
	)

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		c.baseURL+"/v1/oauth2/token",
		strings.NewReader(
			data.Encode(),
		),
	)

	if err != nil {
		return "", err
	}

	req.SetBasicAuth(
		c.clientID,
		c.clientSecret,
	)

	req.Header.Set(
		"Content-Type",
		"application/x-www-form-urlencoded",
	)

	client := &http.Client{}

	resp, err := client.Do(req)

	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {

		return "",
			errors.New(
				"paypal authentication failed",
			)
	}

	var result accessTokenResponse

	if err := json.NewDecoder(
		resp.Body,
	).Decode(&result); err != nil {

		return "", err
	}

	if result.AccessToken == "" {

		return "",
			errors.New(
				"paypal returned empty access token",
			)
	}

	return result.AccessToken, nil
}
