package appErrors

type ValidationErrors map[string]string

type ValidationError struct {
	Errors ValidationErrors `json:"errors"`
}

func NewValidationError(
	errors ValidationErrors,
) *ValidationError {

	return &ValidationError{
		Errors: errors,
	}
}

func (e *ValidationError) Error() string {
	return "validation failed"
}

func (e *ValidationError) HasErrors() bool {
	return len(e.Errors) > 0
}
