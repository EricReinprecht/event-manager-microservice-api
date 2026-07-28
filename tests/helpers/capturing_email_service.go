package helpers

type CapturingEmailService struct {
	VerificationToken string
	ResetToken        string
}

func (c *CapturingEmailService) SendVerificationEmail(
	email string,
	username string,
	token string,
) error {

	c.VerificationToken = token

	return nil
}

func (c *CapturingEmailService) SendPasswordResetEmail(
	email string,
	username string,
	token string,
) error {

	c.ResetToken = token

	return nil
}
