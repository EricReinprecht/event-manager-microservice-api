package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

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
	Username string `json:"username" binding:"required,min=3,max=50"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

type loginRequest struct {
	Identifier string `json:"identifier" binding:"required"`
	Password   string `json:"password" binding:"required"`
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req registerRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
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

		// password validation
		var passwordErr *security.PasswordValidationError
		if errors.As(err, &passwordErr) {
			c.JSON(http.StatusBadRequest, gin.H{
				"field":  "password",
				"errors": passwordErr.Errors,
			})
			return
		}

		// duplicate email
		if err.Error() == "email already exists" {
			c.JSON(http.StatusBadRequest, gin.H{
				"field": "email",
				"error": err.Error(),
			})
			return
		}

		// duplicate username
		if err.Error() == "username already exists" {
			c.JSON(http.StatusBadRequest, gin.H{
				"field": "username",
				"error": err.Error(),
			})
			return
		}

		// fallback
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":    user.ID,
		"email": user.Email,
	})
}

func (h *AuthHandler) Login(c *gin.Context) {

	var req loginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	token, err := h.service.Login(
		c.Request.Context(),
		req.Identifier,
		req.Password,
	)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": token,
	})
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
