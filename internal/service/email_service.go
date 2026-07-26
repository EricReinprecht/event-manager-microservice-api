package service

import (
	"fmt"
)

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
