package auth_service

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
