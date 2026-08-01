package i18n

import (
	"fmt"
	"strings"
)

const DefaultLanguage = "en"

type Registry struct {
	translators map[string]*Translator
}

func NewRegistry(
	languages ...string,
) (*Registry, error) {

	translators := make(
		map[string]*Translator,
		len(languages),
	)

	for _, language := range languages {

		translator, err := NewTranslator(
			language,
		)

		if err != nil {
			return nil, fmt.Errorf(
				"load translations for %s: %w",
				language,
				err,
			)
		}

		translators[language] = translator
	}

	if _, exists := translators[DefaultLanguage]; !exists {
		return nil, fmt.Errorf(
			"default language %q was not loaded",
			DefaultLanguage,
		)
	}

	return &Registry{
		translators: translators,
	}, nil
}

func (r *Registry) Translator(
	language string,
) *Translator {

	language = normalizeLanguage(language)

	if translator, exists := r.translators[language]; exists {
		return translator
	}

	return r.translators[DefaultLanguage]
}

func normalizeLanguage(
	language string,
) string {

	language = strings.TrimSpace(
		strings.ToLower(language),
	)

	if language == "" {
		return DefaultLanguage
	}

	if separator := strings.IndexAny(
		language,
		"-_",
	); separator >= 0 {

		language = language[:separator]
	}

	return language
}
