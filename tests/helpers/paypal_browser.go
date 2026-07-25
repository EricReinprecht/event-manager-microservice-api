package helpers

import (
	"fmt"
	"os"
	"time"

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
			Headless: playwright.Bool(true),
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
		playwright.PageGotoOptions{
			WaitUntil: playwright.WaitUntilStateNetworkidle,
		},
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

	if err := page.
		Locator("#email").
		Fill(email); err != nil {
		return err
	}

	if err := page.
		Locator("#btnNext").
		Click(); err != nil {
		return err
	}

	if err := page.
		Locator("#password").
		Fill(password); err != nil {
		return err
	}

	if err := page.
		Locator("#btnLogin").
		Click(); err != nil {
		return err
	}

	button := page.Locator(
		"button:has-text('Kauf abschließen'), button:has-text('Pay Now')",
	)

	if err := button.WaitFor(); err != nil {
		return err
	}

	if err := button.Click(); err != nil {
		return err
	}

	// IMPORTANT:
	// PayPal approval is asynchronous.
	// Wait until PayPal redirects back to your application.
	err = page.WaitForURL(
		"**?token=*",
		playwright.PageWaitForURLOptions{
			Timeout: playwright.Float(120000),
		},
	)

	if err != nil {
		return fmt.Errorf(
			"paypal approval redirect failed: %w",
			err,
		)
	}

	// Give PayPal a moment to finalize order state
	time.Sleep(
		2 * time.Second,
	)

	return nil
}
