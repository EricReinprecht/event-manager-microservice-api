package requests

type LogoutRequest struct {
	RefreshToken string `json:"refreshToken" binding:"required"`
}
