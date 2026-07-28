package requests

type ForgotPasswordRequest struct {
	Identifier string `json:"identifier" binding:"required"`
}
