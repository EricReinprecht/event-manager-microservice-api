package app

import (
	"github.com/reinp/event-platform/backend/internal/auth"
	"github.com/reinp/event-platform/backend/internal/clock"
	"github.com/reinp/event-platform/backend/internal/config"
	"github.com/reinp/event-platform/backend/internal/mail"
	"github.com/reinp/event-platform/backend/internal/payment/paypal"
	"github.com/reinp/event-platform/backend/internal/security"
)

type Clients struct {
	JWT *auth.JWT

	Mailer *mail.Mailer

	Clock clock.Clock

	PasswordValidator *security.PasswordValidator

	PayPalClient *paypal.Client
}

func NewClients(
	cfg *config.Config,
) *Clients {

	appClock := clock.RealClock{}

	jwt := auth.NewJWT(
		cfg.JWTSecret,
		appClock,
		cfg.AccessTokenDuration,
	)

	mailer := mail.NewMailer(

		cfg.SMTPHost,

		cfg.SMTPPort,

		cfg.SMTPUser,

		cfg.SMTPPassword,

		cfg.SMTPFrom,
	)

	passwordValidator :=
		security.NewPasswordValidator()

	paypalClient :=
		paypal.NewClient(

			cfg.PayPalClientID,

			cfg.PayPalClientSecret,

			cfg.PayPalBaseURL,

			cfg.PayPalReturnURL,

			cfg.PayPalCancelURL,

			cfg.PayPalWebhookID,
		)

	return &Clients{

		JWT: jwt,

		Mailer: mailer,

		Clock: appClock,

		PasswordValidator: passwordValidator,

		PayPalClient: paypalClient,
	}
}
