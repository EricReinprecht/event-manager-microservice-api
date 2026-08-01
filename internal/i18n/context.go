package i18n

import (
	"github.com/gin-gonic/gin"
)

const ContextKey = "translator"

func FromContext(
	c *gin.Context,
) *Translator {

	value, exists := c.Get(
		ContextKey,
	)

	if exists {

		if translator, ok := value.(*Translator); ok {
			return translator
		}
	}

	return &Translator{
		translations: map[string]string{},
	}
}
