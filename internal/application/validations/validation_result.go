package validations

type ValidationResult struct {
	Passed bool
	Errors []ValidationError
}
