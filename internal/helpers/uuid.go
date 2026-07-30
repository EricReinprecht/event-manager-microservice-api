package helpers

import "github.com/google/uuid"

func ParseUUIDs(values []string) ([]uuid.UUID, error) {
	result := make([]uuid.UUID, 0, len(values))

	for _, value := range values {
		id, err := uuid.Parse(value)

		if err != nil {
			return nil, err
		}

		result = append(result, id)
	}

	return result, nil
}
