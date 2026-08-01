package i18n

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Translator struct {
	translations map[string]string
}

func NewTranslator(
	language string,
) (*Translator, error) {

	path := filepath.Join(
		"internal",
		"i18n",
		language,
	)

	files, err := filepath.Glob(
		filepath.Join(path, "*.json"),
	)

	if err != nil {
		return nil, err
	}

	translations := make(map[string]string)

	for _, file := range files {

		content, err := os.ReadFile(
			file,
		)

		if err != nil {
			return nil, err
		}

		var data map[string]any

		if err := json.Unmarshal(
			content,
			&data,
		); err != nil {
			return nil, err
		}

		flatten(
			"",
			data,
			translations,
		)
	}

	return &Translator{
		translations: translations,
	}, nil
}

func (t *Translator) T(
	key string,
) string {

	if value, ok := t.translations[key]; ok {
		return value
	}

	return key
}

func flatten(
	prefix string,
	data map[string]any,
	target map[string]string,
) {

	for key, value := range data {

		fullKey := key

		if prefix != "" {
			fullKey = prefix + "." + key
		}

		switch v := value.(type) {

		case string:
			target[fullKey] = v

		case map[string]any:
			flatten(
				fullKey,
				v,
				target,
			)

		default:
			continue
		}
	}
}
