package validations

type Validatable interface {
	Validate() ValidationResult
}

type Validator interface {
	Validate(item Validatable) ValidationResult
}
