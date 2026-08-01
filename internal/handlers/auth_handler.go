package handlers

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	appErrors "github.com/reinp/event-platform/backend/internal/appErrors"
	"github.com/reinp/event-platform/backend/internal/helpers"
	"github.com/reinp/event-platform/backend/internal/requests"
	"github.com/reinp/event-platform/backend/internal/responses"
	"github.com/reinp/event-platform/backend/internal/security"
	"github.com/reinp/event-platform/backend/internal/service"
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

func (h *AuthHandler) Register(c *gin.Context) {

	var req requests.RegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {

		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": security.ErrorMessage(err),
			},
		)

		return
	}

	user, err := h.service.Register(
		c.Request.Context(),
		service.RegisterRequest{
			Email:    req.Email,
			Password: req.Password,
			Username: req.Username,
		},
	)

	if err != nil {

		// Password validation
		if passwordErr, ok := errors.AsType[*security.PasswordValidationError](err); ok {

			c.JSON(
				http.StatusBadRequest,
				gin.H{
					"field":  "password",
					"errors": passwordErr.Errors,
				},
			)

			return
		}

		// Username already exists
		// Safe to expose because usernames are public identifiers

		if errors.Is(
			err,
			appErrors.ErrUsernameAlreadyExists,
		) {

			c.JSON(
				http.StatusConflict,
				gin.H{
					"field": "username",
					"error": "username already exists",
				},
			)

			return
		}

		// Email already exists
		// Do NOT reveal account existence

		if errors.Is(
			err,
			appErrors.ErrEmailAlreadyExists,
		) {

			c.JSON(
				http.StatusCreated,
				gin.H{
					"message": "Please check your email to complete registration.",
				},
			)

			return
		}

		// fallback

		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": err.Error(),
			},
		)

		return
	}

	c.JSON(
		http.StatusCreated,
		gin.H{
			"message": "Please check your email to complete registration.",
			"id":      user.ID,
		},
	)
}

func (h *AuthHandler) Login(c *gin.Context) {

	var req requests.LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {

		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": err.Error(),
			},
		)

		return
	}

	response, err := h.service.Login(
		c.Request.Context(),
		req.Identifier,
		req.Password,
		c.Request.UserAgent(),
		c.ClientIP(),
	)

	if err != nil {
		responses.HandleDomainError(c, err)
		return
	}

	helpers.SetRefreshTokenCookie(
		c,
		response.RefreshToken,
		h.refreshTokenDuration,
		h.cookieSecure,
	)

	responses.Success(
		c,
		http.StatusOK,
		gin.H{
			"accessToken": response.AccessToken,
		},
	)
}

func (h *AuthHandler) VerifyEmail(c *gin.Context) {

	token := c.Query("token")

	if token == "" {

		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": "missing token",
			},
		)

		return
	}

	jwtToken, err := h.service.VerifyEmail(
		c.Request.Context(),
		token,
	)

	if err != nil {

		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": err.Error(),
			},
		)

		return
	}

	c.JSON(
		http.StatusOK,
		gin.H{
			"token": jwtToken,
		},
	)
}

func (h *AuthHandler) Refresh(
	c *gin.Context,
) {
	refreshToken, err := c.Cookie(
		helpers.RefreshTokenCookie,
	)

	if err != nil || refreshToken == "" {
		responses.Unauthorized(c)
		return
	}

	response, err := h.service.Refresh(
		c.Request.Context(),
		refreshToken,
	)

	if err != nil {
		helpers.ClearRefreshTokenCookie(
			c,
			h.cookieSecure,
		)

		responses.Unauthorized(c)
		return
	}

	// Refresh-token rotation:
	// replace the old cookie with the newly generated token.
	helpers.SetRefreshTokenCookie(
		c,
		response.RefreshToken,
		h.refreshTokenDuration,
		h.cookieSecure,
	)

	responses.Success(
		c,
		http.StatusOK,
		gin.H{
			"accessToken": response.AccessToken,
		},
	)
}

func (h *AuthHandler) Logout(
	c *gin.Context,
) {
	refreshToken, err := c.Cookie(
		helpers.RefreshTokenCookie,
	)

	if err == nil && refreshToken != "" {
		_ = h.service.Logout(
			c.Request.Context(),
			refreshToken,
		)
	}

	helpers.ClearRefreshTokenCookie(
		c,
		h.cookieSecure,
	)

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

		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": err.Error(),
			},
		)

		return
	}

	err := h.service.ForgotPassword(
		c.Request.Context(),
		req.Identifier,
	)

	// Always return success to prevent email enumeration
	if err != nil {

		c.JSON(
			http.StatusOK,
			gin.H{
				"message": "if this email exists, a reset link was sent",
			},
		)

		return
	}

	c.JSON(
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

		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": err.Error(),
			},
		)

		return
	}

	err := h.service.ResetPassword(
		c.Request.Context(),
		req.Token,
		req.NewPassword,
	)

	if err != nil {

		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": err.Error(),
			},
		)

		return
	}

	c.JSON(
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

		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": err.Error(),
			},
		)

		return
	}

	err := h.service.ResendVerificationEmail(
		c.Request.Context(),
		req.Email,
	)

	if err != nil {

		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": err.Error(),
			},
		)

		return
	}

	c.JSON(
		http.StatusOK,
		gin.H{
			"message": "verification email sent",
		},
	)
}
