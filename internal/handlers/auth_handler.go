package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/reinp/event-platform/backend/internal/appErrors"
	"github.com/reinp/event-platform/backend/internal/helpers"
	"github.com/reinp/event-platform/backend/internal/requests"
	"github.com/reinp/event-platform/backend/internal/responses"
	"github.com/reinp/event-platform/backend/internal/service"
	auth_service "github.com/reinp/event-platform/backend/internal/service/auth"
)

const (
	registrationSuccessMessage = "Please check your email to complete registration."

	passwordResetRequestMessage = "if this email exists, a reset link was sent"
)

type AuthHandler struct {
	service              *service.AuthService
	refreshTokenDuration time.Duration
	cookieSecure         bool
}

func NewAuthHandler(
	service *service.AuthService,
	refreshTokenDuration time.Duration,
	cookieSecure bool,
) *AuthHandler {

	return &AuthHandler{
		service:              service,
		refreshTokenDuration: refreshTokenDuration,
		cookieSecure:         cookieSecure,
	}
}

func (h *AuthHandler) Register(
	c *gin.Context,
) {

	var req requests.RegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {

		responses.BadRequest(
			c,
			err,
		)

		return
	}

	user, err := h.service.Register(
		c.Request.Context(),
		auth_service.RegisterRequest{
			Email:    req.Email,
			Password: req.Password,
			Username: req.Username,
		},
	)

	if err != nil {

		responses.HandleDomainError(
			c,
			err,
		)

		return
	}

	responses.Success(
		c,
		http.StatusCreated,
		gin.H{
			"message": "Please check your email to complete registration.",
			"id":      user.ID,
		},
	)
}

func (h *AuthHandler) Login(
	c *gin.Context,
) {

	var req requests.LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {

		responses.BadRequest(
			c,
			err,
		)

		return
	}

	result, err := h.service.Login(
		c.Request.Context(),
		req.Identifier,
		req.Password,
		c.Request.UserAgent(),
		c.ClientIP(),
	)

	if err != nil {

		responses.HandleDomainError(
			c,
			err,
		)

		return
	}

	h.setRefreshTokenCookie(
		c,
		result.RefreshToken,
	)

	responses.Success(
		c,
		http.StatusOK,
		gin.H{
			"accessToken": result.AccessToken,
		},
	)
}

func (h *AuthHandler) VerifyEmail(
	c *gin.Context,
) {

	token := c.Query("token")

	if token == "" {

		responses.BadRequest(
			c,
			appErrors.ErrMissingToken,
		)

		return
	}

	accessToken, err := h.service.VerifyEmail(
		c.Request.Context(),
		token,
	)

	if err != nil {

		responses.HandleDomainError(
			c,
			err,
		)

		return
	}

	responses.Success(
		c,
		http.StatusOK,
		gin.H{
			"accessToken": accessToken,
		},
	)
}

func (h *AuthHandler) Refresh(
	c *gin.Context,
) {

	refreshToken, err := h.readRefreshTokenCookie(c)

	if err != nil {

		responses.Unauthorized(c)

		return
	}

	result, err := h.service.Refresh(
		c.Request.Context(),
		refreshToken,
	)

	if err != nil {

		h.clearRefreshTokenCookie(c)

		responses.Unauthorized(c)

		return
	}

	h.setRefreshTokenCookie(
		c,
		result.RefreshToken,
	)

	responses.Success(
		c,
		http.StatusOK,
		gin.H{
			"accessToken": result.AccessToken,
		},
	)
}

func (h *AuthHandler) Logout(
	c *gin.Context,
) {

	refreshToken, _ := h.readRefreshTokenCookie(c)

	if refreshToken != "" {

		_ = h.service.Logout(
			c.Request.Context(),
			refreshToken,
		)
	}

	h.clearRefreshTokenCookie(c)

	responses.Success(
		c,
		http.StatusOK,
		gin.H{
			"message": "logged out",
		},
	)
}

func (h *AuthHandler) ForgotPassword(
	c *gin.Context,
) {

	var req requests.ForgotPasswordRequest

	if err := c.ShouldBindJSON(&req); err != nil {

		responses.BadRequest(
			c,
			err,
		)

		return
	}

	err := h.service.ForgotPassword(
		c.Request.Context(),
		req.Identifier,
	)

	if err != nil {

		responses.HandleDomainError(
			c,
			err,
		)

		return
	}

	responses.Success(
		c,
		http.StatusOK,
		gin.H{
			"message": "if this email exists, a reset link was sent",
		},
	)
}

func (h *AuthHandler) ResetPassword(
	c *gin.Context,
) {

	var req requests.ResetPasswordRequest

	if err := c.ShouldBindJSON(&req); err != nil {

		responses.BadRequest(
			c,
			err,
		)

		return
	}

	if err := h.service.ResetPassword(
		c.Request.Context(),
		req.Token,
		req.NewPassword,
	); err != nil {

		responses.HandleDomainError(
			c,
			err,
		)

		return
	}

	h.clearRefreshTokenCookie(c)

	responses.Success(
		c,
		http.StatusOK,
		gin.H{
			"message": "password reset successful",
		},
	)
}

func (h *AuthHandler) ResendVerificationEmail(
	c *gin.Context,
) {

	var req requests.ResendVerificationRequest

	if err := c.ShouldBindJSON(&req); err != nil {

		responses.BadRequest(
			c,
			err,
		)

		return
	}

	if err := h.service.ResendVerificationEmail(
		c.Request.Context(),
		req.Email,
	); err != nil {

		responses.HandleDomainError(
			c,
			err,
		)

		return
	}

	responses.Success(
		c,
		http.StatusOK,
		gin.H{
			"message": "verification email sent",
		},
	)
}

func (h *AuthHandler) readRefreshTokenCookie(
	c *gin.Context,
) (string, error) {

	refreshToken, err := c.Cookie(
		helpers.RefreshTokenCookie,
	)

	if err != nil || refreshToken == "" {
		return "", http.ErrNoCookie
	}

	return refreshToken, nil
}

func (h *AuthHandler) setRefreshTokenCookie(
	c *gin.Context,
	refreshToken string,
) {

	helpers.SetRefreshTokenCookie(
		c,
		refreshToken,
		h.refreshTokenDuration,
		h.cookieSecure,
	)
}

func (h *AuthHandler) clearRefreshTokenCookie(
	c *gin.Context,
) {

	helpers.ClearRefreshTokenCookie(
		c,
		h.cookieSecure,
	)
}
