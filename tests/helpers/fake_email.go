package helpers

import "errors"

type FailingEmailService struct{}

func (f *FailingEmailService) SendVerificationEmail(
	email string,
	username string,
	token string,
) error {

	return errors.New(
		"mail delivery failed",
	)
}

func (f *FailingEmailService) SendPasswordResetEmail(
	email string,
	username string,
	token string,
) error {

	return errors.New(
		"mail delivery failed",
	)
}
