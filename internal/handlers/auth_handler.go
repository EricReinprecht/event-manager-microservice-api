package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	appErrors "github.com/reinp/event-platform/backend/internal/appErrors"
	"github.com/reinp/event-platform/backend/internal/security"
	"github.com/reinp/event-platform/backend/internal/service"
)

type AuthHandler struct {
	service *service.AuthService
}

func NewAuthHandler(
	service *service.AuthService,
) *AuthHandler {

	return &AuthHandler{
		service: service,
	}
}

type registerRequest struct {
	Username string `json:"username" binding:"required,min=3,max=30"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=12,max=128"`
}

type loginRequest struct {
	Identifier string `json:"identifier" binding:"required"`
	Password   string `json:"password" binding:"required"`
}

type logoutRequest struct {
	RefreshToken string `json:"refreshToken" binding:"required"`
}

type ForgotPasswordRequest struct {
	Identifier string `json:"identifier" binding:"required"`
}

type ResetPasswordRequest struct {
	Token       string `json:"token" binding:"required"`
	NewPassword string `json:"newPassword" binding:"required"`
}

func (h *AuthHandler) Register(c *gin.Context) {

	var req registerRequest

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

	var req loginRequest

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

		c.JSON(
			http.StatusUnauthorized,
			gin.H{
				"error": err.Error(),
			},
		)

		return
	}

	c.JSON(
		http.StatusOK,
		gin.H{
			"accessToken":  response.AccessToken,
			"refreshToken": response.RefreshToken,
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

	var req struct {
		RefreshToken string `json:"refreshToken" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {

		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": "invalid request",
			},
		)

		return
	}

	response, err := h.service.Refresh(
		c.Request.Context(),
		req.RefreshToken,
	)

	if err != nil {

		c.JSON(
			http.StatusUnauthorized,
			gin.H{
				"error": err.Error(),
			},
		)

		return
	}

	c.JSON(
		http.StatusOK,
		gin.H{
			"accessToken":  response.AccessToken,
			"refreshToken": response.RefreshToken,
		},
	)
}

func (h *AuthHandler) Logout(
	c *gin.Context,
) {

	var req logoutRequest

	if err := c.ShouldBindJSON(&req); err != nil {

		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": err.Error(),
			},
		)

		return
	}

	err := h.service.Logout(
		c.Request.Context(),
		req.RefreshToken,
	)

	if err != nil {

		c.JSON(
			http.StatusUnauthorized,
			gin.H{
				"error": "invalid refresh token",
			},
		)

		return
	}

	c.JSON(
		http.StatusOK,
		gin.H{
			"message": "logged out successfully",
		},
	)
}

func (h *AuthHandler) ForgotPassword(
	c *gin.Context,
) {

	var req ForgotPasswordRequest

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

	var req ResetPasswordRequest

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
