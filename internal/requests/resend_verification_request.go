package requests

type ResendVerificationRequest struct {
	Email string `json:"email" binding:"required,email"`
}
