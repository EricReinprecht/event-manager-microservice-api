package helpers

import (
	"os"

	"github.com/mxschmitt/playwright-go"
)

func ApprovePayPalOrder(
	url string,
) error {

	pw, err := playwright.Run()

	if err != nil {
		return err
	}

	defer pw.Stop()

	browser, err := pw.Chromium.Launch(
		playwright.BrowserTypeLaunchOptions{
			Headless: playwright.Bool(false),
		},
	)

	if err != nil {
		return err
	}

	defer browser.Close()

	page, err := browser.NewPage()

	if err != nil {
		return err
	}

	page.SetDefaultTimeout(120000)

	_, err = page.Goto(
		url,
	)

	if err != nil {
		return err
	}

	email := os.Getenv(
		"PAYPAL_BUYER_EMAIL",
	)

	password := os.Getenv(
		"PAYPAL_BUYER_PASSWORD",
	)

	err = page.
		Locator("#email").
		Fill(email)

	if err != nil {
		return err
	}

	err = page.
		Locator("#btnNext").
		Click()

	if err != nil {
		return err
	}

	err = page.
		Locator("#password").
		Fill(password)

	if err != nil {
		return err
	}

	err = page.
		Locator("#btnLogin").
		Click()

	if err != nil {
		return err
	}

	button := page.Locator(
		"button:has-text('Kauf abschließen'), button:has-text('Pay Now')",
	)

	err = button.WaitFor()

	if err != nil {
		return err
	}

	err = button.Click()

	if err != nil {
		return err
	}

	return nil
}
