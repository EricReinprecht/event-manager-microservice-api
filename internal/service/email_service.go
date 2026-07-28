package service

import (
	"fmt"
)

type EmailSender interface {
	SendVerificationEmail(
		email string,
		username string,
		token string,
	) error

	SendPasswordResetEmail(
		email string,
		username string,
		token string,
	) error
}

type EmailService struct {
	mailer      MailSender
	frontendURL string
}

type MailSender interface {
	Send(
		to string,
		subject string,
		body string,
	) error
}

func NewEmailService(
	mailer MailSender,
	frontendURL string,
) *EmailService {

	return &EmailService{
		mailer:      mailer,
		frontendURL: frontendURL,
	}
}

func (s *EmailService) SendVerificationEmail(
	email string,
	username string,
	token string,
) error {

	link := fmt.Sprintf(
		"%s/verify-email?token=%s",
		s.frontendURL,
		token,
	)

	body := fmt.Sprintf(
		`
Hello %s,

Welcome to Event Platform.

Please verify your email:

%s

This link expires in 24 hours.
`,
		username,
		link,
	)

	return s.mailer.Send(
		email,
		"Verify your email",
		body,
	)
}

func (s *EmailService) SendPasswordResetEmail(
	email string,
	username string,
	token string,
) error {

	resetURL := fmt.Sprintf(
		"%s/reset-password?token=%s",
		s.frontendURL,
		token,
	)

	subject := "Reset your password"

	body := fmt.Sprintf(
		`
Hello %s,

You requested a password reset.

Click the link below to reset your password:

%s

This link expires in 15 minutes.

If you did not request this, ignore this email.

Regards
Event Platform
`,
		username,
		resetURL,
	)

	return s.mailer.Send(
		email,
		subject,
		body,
	)
}
