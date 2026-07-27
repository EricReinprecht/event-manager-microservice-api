package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"

	"github.com/reinp/event-platform/backend/internal/service"
)

func Auth(
	authService *service.AuthService,
) gin.HandlerFunc {

	return func(c *gin.Context) {

		header := c.GetHeader("Authorization")

		if header == "" {

			c.JSON(
				http.StatusUnauthorized,
				gin.H{
					"error": "missing authorization header",
				},
			)

			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(
			header,
			"Bearer ",
		)

		token, err := jwt.Parse(
			tokenString,
			func(token *jwt.Token) (interface{}, error) {

				return []byte(authService.Secret()), nil
			},
		)

		if err != nil || !token.Valid {

			c.JSON(
				http.StatusUnauthorized,
				gin.H{
					"error": "invalid token",
				},
			)

			c.Abort()
			return
		}

		claims := token.Claims.(jwt.MapClaims)

		userID, ok := claims["user_id"].(string)

		if !ok {

			c.JSON(
				http.StatusUnauthorized,
				gin.H{
					"error": "invalid user id",
				},
			)

			c.Abort()
			return
		}

		userUUID, err := uuid.Parse(userID)

		if err != nil {

			c.JSON(
				http.StatusUnauthorized,
				gin.H{
					"error": "invalid user id",
				},
			)

			c.Abort()
			return
		}

		user, err := authService.ValidateUser(
			c.Request.Context(),
			userUUID,
		)

		if err != nil {

			c.JSON(
				http.StatusUnauthorized,
				gin.H{
					"error": "user no longer exists",
				},
			)

			c.Abort()
			return
		}

		c.Set(
			"user_id",
			user.ID,
		)

		c.Set(
			"user",
			user,
		)

		c.Next()
	}
}
