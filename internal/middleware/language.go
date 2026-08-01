package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/reinp/event-platform/backend/internal/i18n"
)

const TranslatorContextKey = "translator"

func Language(
	registry *i18n.Registry,
) gin.HandlerFunc {

	return func(c *gin.Context) {

		language := readLanguage(c)

		translator := registry.Translator(
			language,
		)

		c.Set(
			TranslatorContextKey,
			translator,
		)

		c.Next()
	}
}

func readLanguage(
	c *gin.Context,
) string {

	// Optional explicit query override:
	// /api/parties?lang=de
	if language := c.Query("lang"); language != "" {
		return language
	}

	header := c.GetHeader(
		"Accept-Language",
	)

	if header == "" {
		return i18n.DefaultLanguage
	}

	firstLanguage := strings.Split(
		header,
		",",
	)[0]

	firstLanguage = strings.Split(
		firstLanguage,
		";",
	)[0]

	return strings.TrimSpace(
		firstLanguage,
	)
}
