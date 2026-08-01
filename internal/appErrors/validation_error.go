package appErrors

const (
	ErrMsgPartyEndBeforeStart = "party.end_before_start"

	ErrMsgLocationRequired = "party.location_required"

	ErrMsgTimezoneRequired = "party.timezone_required"
)

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
