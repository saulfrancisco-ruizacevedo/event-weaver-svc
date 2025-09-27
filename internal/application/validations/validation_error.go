package validations

import "fmt"

type ValidationError struct {
	EntityName   string
	EntityType   string // e.g., "event", "component", "team", "domain"
	ErrorCode    string // A code to classify the error (e.g., "NON_EXISTENT_REFERENCE")
	ErrorMessage string // A human-readable message
}

func (e ValidationError) Error() string {
	return fmt.Sprintf("validation error in %s '%s': %s", e.EntityType, e.EntityName, e.ErrorMessage)
}

const (
	NonExistentReference = "NON_EXISTENT_REFERENCE"
	EmptyField           = "EMPTY_FIELD"
)
