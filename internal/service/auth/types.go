package auth_service

import (
	"time"

	"github.com/google/uuid"
)

type RegisterResult struct {
	UserID *uuid.UUID
}

type RegisterRequest struct {
	Email    string
	Password string
	Username string
}

type TokenResponse struct {
	AccessToken  string
	RefreshToken string
}

type SessionResponse struct {
	FamilyID  uuid.UUID `json:"familyId"`
	Device    string    `json:"device"`
	IP        string    `json:"ip"`
	CreatedAt time.Time `json:"createdAt"`
	Current   bool      `json:"current"`
}

type ForgotPasswordRequest struct {
	Identifier string
}

type ResetPasswordRequest struct {
	Token       string
	NewPassword string
}
