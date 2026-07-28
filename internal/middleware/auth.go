package middleware

import (
	"errors"
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

		header := c.GetHeader(
			"Authorization",
		)

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

		parts := strings.Split(
			header,
			" ",
		)

		if len(parts) != 2 || parts[0] != "Bearer" {

			c.JSON(
				http.StatusUnauthorized,
				gin.H{
					"error": "invalid authorization header",
				},
			)

			c.Abort()
			return
		}

		tokenString := parts[1]

		token, err := jwt.Parse(
			tokenString,
			func(token *jwt.Token) (interface{}, error) {

				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {

					return nil, errors.New(
						"unexpected signing method",
					)
				}

				return []byte(
					authService.Secret(),
				), nil
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

		claims, ok := token.Claims.(jwt.MapClaims)

		if !ok {

			c.JSON(
				http.StatusUnauthorized,
				gin.H{
					"error": "invalid claims",
				},
			)

			c.Abort()
			return
		}

		// ======================
		// USER ID
		// ======================

		userIDString, ok := claims["user_id"].(string)

		if !ok {

			c.JSON(
				http.StatusUnauthorized,
				gin.H{
					"error": "missing user id",
				},
			)

			c.Abort()
			return
		}

		userUUID, err := uuid.Parse(
			userIDString,
		)

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

		// ======================
		// FAMILY ID
		// ======================

		familyIDString, ok := claims["family_id"].(string)

		if !ok {

			c.JSON(
				http.StatusUnauthorized,
				gin.H{
					"error": "missing family id",
				},
			)

			c.Abort()
			return
		}

		familyUUID, err := uuid.Parse(
			familyIDString,
		)

		if err != nil {

			c.JSON(
				http.StatusUnauthorized,
				gin.H{
					"error": "invalid family id",
				},
			)

			c.Abort()
			return
		}

		// ======================
		// USER EXISTS
		// ======================

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

		// ======================
		// CONTEXT
		// ======================

		c.Set(
			"userID",
			user.ID,
		)

		c.Set(
			"familyID",
			familyUUID,
		)

		c.Set(
			"user",
			user,
		)

		c.Next()
	}
}
