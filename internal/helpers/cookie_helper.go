package helpers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

const RefreshTokenCookie = "refresh_token"

func SetRefreshTokenCookie(
	c *gin.Context,
	token string,
	duration time.Duration,
	secure bool,
) {
	http.SetCookie(
		c.Writer,
		&http.Cookie{
			Name:     RefreshTokenCookie,
			Value:    token,
			Path:     "/api/auth",
			MaxAge:   int(duration.Seconds()),
			Expires:  time.Now().Add(duration),
			HttpOnly: true,
			Secure:   secure,
			SameSite: http.SameSiteLaxMode,
		},
	)
}

func ClearRefreshTokenCookie(
	c *gin.Context,
	secure bool,
) {
	http.SetCookie(
		c.Writer,
		&http.Cookie{
			Name:     RefreshTokenCookie,
			Value:    "",
			Path:     "/api/auth",
			MaxAge:   -1,
			Expires:  time.Unix(0, 0),
			HttpOnly: true,
			Secure:   secure,
			SameSite: http.SameSiteLaxMode,
		},
	)
}
